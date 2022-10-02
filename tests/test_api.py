from conftest import PUBLIC, INTERNAL


def test_public_hello(ctx):
    """
    Check public hello
    """
    res = ctx.reporting_get(PUBLIC, "/api/v1/hello")
    assert res.status_code == 200
    assert res.json() == {"message": "Hello world public"}


def test_internal_hello(ctx):
    """
    Check internal hello
    """
    res = ctx.reporting_get(INTERNAL, "/api/v1/hello")
    assert res.status_code == 200
    assert res.json() == {"message": "Hello world internal"}
