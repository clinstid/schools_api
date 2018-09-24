from http import HTTPStatus

from common import (
    TestSchoolsAPI,
    check_error_response,
)


class TestGetSchool(TestSchoolsAPI):
    def test_get_school_simple(self):
        response = self.get_school(school_id=0)
        assert response.status_code == HTTPStatus.OK
        school = response.json()
        for field in ('name', 'id'):
            assert field in school

    def test_get_school_not_found(self):
        bad_id = 1000000000
        response = self.get_school(school_id=bad_id)
        assert response.status_code == HTTPStatus.NOT_FOUND
        check_error_response(response, f'School with id {bad_id} not found')

    def test_get_school_invalid_id(self):
        bad_id = 'notanumber'
        response = self.get_school(school_id=bad_id)
        assert response.status_code == HTTPStatus.BAD_REQUEST
        check_error_response(response, 'school id must be a number')
