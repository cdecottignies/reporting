from pathlib import Path

import docker
import pytest
import requests
from pyvade.test.context import Context
from pyvade.test.utils import Wait, archive2file

ignore = type("ignore", (), {"__eq__": lambda x, y: True})()

PUBLIC = 8080
INTERNAL = 8081
MONITORING = 8082


class Reporting(Context):
    def __init__(self, **kwargs):
        Context.__init__(self, **kwargs)

    def reporting_base_addr(self, port):
        return "http://%s:%d" % (self.container_addr("reporting"), port)

    def wait_for_running_state(self):
        Wait(
            timeout=60,
            ignored_exceptions=(requests.exceptions.ConnectionError,),
            cb_check_result=(lambda x: x.status_code == 200),
        ).until(
            lambda: requests.get(self.reporting_base_addr(MONITORING) + "/info/ready"),
            "reporting is ready?",
        )

    def request(self, method, url, **kwargs):
        try:
            r = requests.request(method, url, **kwargs)
            dump_req_resp(r)
            return r
        except requests.RequestException as e:
            dump_req_resp_from_exception(e)
            raise

    def reporting_request(self, method, port, path, role=None, **kwargs):
        return self.request(
            method,
            "%s%s" % (self.reporting_base_addr(port), path),
            role=role,
            **kwargs,
        )

    def reporting_get(self, port, path, params=None, **kwargs):
        return self.request(
            "GET",
            "%s%s" % (self.reporting_base_addr(port), path),
            params=params,
            **kwargs,
        )

    def reporting_put(self, port, path, json=None, **kwargs):
        return self.request(
            "PUT",
            "%s%s" % (self.reporting_base_addr(port), path),
            json=json,
            **kwargs,
        )


def dump_req_resp(r):
    print("---------------- REQUEST ----------------------")
    print(r.request.method + " " + r.request.url)
    print("Headers: " + str(r.request.headers))
    if r.request.body is not None:
        print(r.request.body)
    print("---------------- RESPONSE ---------------------")
    print("%d %s" % (r.status_code, r.reason))
    print("Headers: " + str(r.headers))
    if r.text != "":
        print(r.text)
    print("---------------- END ---------------------")


def dump_req_resp_from_exception(e):
    print("---------------- REQUEST (exception) ----------------------")
    print(e.request.method + " " + e.request.url)
    print("Headers: " + str(e.request.headers))
    if e.request.body is not None:
        print(e.request.body)
    print("---------------- RESPONSE (exception) ---------------------")
    if e.response is None:
        print("NO RESPONSE RECEIVED")
    else:
        print("%d %s" % (e.response.status_code, e.response.reason))
        print("Headers: " + str(e.response.headers))
        if e.response.text != "":
            print(e.response.text)
    print("---------------- END ---------------------")


@pytest.fixture(scope="session")
def platform(dockerc, dockerc_logs):
    ptf = Reporting(docker_compose=dockerc)
    ptf.wait_for_running_state()
    yield ptf


@pytest.fixture
def ctx(platform):
    yield platform


def pytest_sessionfinish(session, exitstatus):
    if exitstatus != 0:
        return

    print("Retrieving coverage data from functional tests...")

    covpath = Path(session.config.rootdir) / "build" / "coverage"
    covpath.mkdir(parents=True, exist_ok=True)
    covfile = covpath / "functional_test.cov"

    client = docker.from_env()
    container = client.containers.run(
        "alpine",
        auto_remove=True,
        detach=True,
        tty=True,
        volumes={"reporting-coverage": {"bind": "/coverage", "mode": "rw"}},
    )
    try:
        report, _ = container.get_archive(Path("/coverage/functional_tests.cov"))
    except docker.errors.APIError:
        report = None
    finally:
        container.stop()

    if report:
        with covfile.open("w") as fd:
            print("Coverage file retrieved")
            fd.write(archive2file(report).decode())
