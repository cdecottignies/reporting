import requests

from conftest import MONITORING, PUBLIC, ignore


def test_public_ping(ctx):
    """
    Public ping must return 204 OK
    """
    res = ctx.reporting_get(PUBLIC, "/api/v1/ping")
    assert res.status_code == 204


def test_info_ready(ctx):
    """
    Check the ready endpoint returns 200 when the moderator service is ready to accept requests
    """
    res = ctx.reporting_get(MONITORING, "/info/ready")
    assert res.status_code == requests.codes.ok


def test_info_alive(ctx):
    """
    Check the alive endpoint returns 200 when the moderator service is up
    """
    res = ctx.reporting_get(MONITORING, "/info/alive")
    assert res.status_code == requests.codes.ok


def test_info_operating(ctx):
    """
    Check the operating endpoint returns 200 OK
    """
    res = ctx.reporting_get(MONITORING, "/info/operating")
    assert res.status_code == requests.codes.ok
    assert res.text == "OK"


def test_info_version(ctx):
    """
    Check the version endpoint returns 200 OK with version
    """
    res = ctx.reporting_get(MONITORING, "/info/version")
    assert res.status_code == requests.codes.ok
    assert res.json() == {"version": ignore}


def test_metrics(ctx):
    """
    Check that the returned metrics contains some of the expected metrics.
    This test is pretty basic and doesn't check the returned metrics in depth.
    """
    res = ctx.reporting_get(MONITORING, "/info/metrics")
    assert res.status_code == requests.codes.ok
    assert "go_threads" in res.text


def test_log(ctx):
    """
    Check we can set the log level dynamically and retrieve it
    """
    res = ctx.reporting_get(MONITORING, "/log/level")
    assert res.status_code == 200
    assert res.json() == {
        "level": "info",
        "available_levels": ["debug", "info", "warn", "error"],
    }
    default_level = res.json()["level"]

    # invalid log level
    res = ctx.reporting_put(MONITORING, "/log/level/unknown")
    assert res.status_code == 400

    res = ctx.reporting_get(MONITORING, "/log/level")
    assert res.status_code == 200
    assert res.json()["level"] == default_level

    for level in ["deBug", "info", "warn", "error"]:
        res = ctx.reporting_put(MONITORING, f"/log/level/{level}")
        assert res.status_code == 204

        res = ctx.reporting_get(MONITORING, "/log/level")
        assert res.status_code == 200
        assert res.json()["level"] == level.lower()

    # reset default log level
    res = ctx.reporting_put(MONITORING, f"/log/level/{default_level}")
    assert res.status_code == 204
