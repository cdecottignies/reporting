import pathlib
import subprocess

import pytest

from conftest import PUBLIC


@pytest.mark.parametrize(
    "test_name",
    [
        "TestFunctionalPing",
        "TestFunctionalHello",
    ],
)
def _test_http_go_client_unit_tests(ctx, test_name):
    """
    This test runs the test binary of the client package with the real reporting service.
    """
    build_path = pathlib.Path(__file__).resolve().parents[1] / "build"

    client_test_bin = str(build_path / "client" / "reporting.test")
    coverage_path = str(
        build_path / "coverage" / ("unit_tests_go_client_%s.cov" % test_name)
    )

    subprocess.run(
        [
            client_test_bin,
            "-test.run",
            test_name,
            "-test.v",
            "-test.coverprofile",
            coverage_path,
        ],
        check=True,
        env={
            "TESTS_REPORTING_BASE_ADDR": ctx.reporting_base_addr(PUBLIC),
        },
    )
