from http import HTTPStatus

from common import (
    TestSchoolsAPI,
    check_error_response,
)


class TestUpdateSchool(TestSchoolsAPI):
    def test_update_school_simple(self):
        new_name = 'New Name'
        response = self.update_school(school_id=0, name=new_name)
        assert response.status_code == HTTPStatus.OK
        school = response.json()
        for field in ('name', 'id'):
            assert field in school

        assert school.get('name') == new_name

    def test_update_school_not_found(self):
        bad_id = 1000000000
        response = self.update_school(school_id=bad_id, name="bad id name")
        assert response.status_code == HTTPStatus.NOT_FOUND
        check_error_response(response, f'School with id {bad_id} not found')

    def test_update_school_invalid_id(self):
        bad_id = 'notanumber'
        response = self.update_school(school_id=bad_id, name="invalid id name")
        assert response.status_code == HTTPStatus.BAD_REQUEST
        check_error_response(response, 'school id must be a number')
