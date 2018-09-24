from http import HTTPStatus

from common import (
    TestSchoolsAPI,
    check_error_response,
)


class TestAddSchool(TestSchoolsAPI):
    def test_add_school_simple(self):
        new_name = 'New School Name'
        response = self.add_school(name=new_name)
        assert response.status_code == HTTPStatus.CREATED

        school = response.json()
        for field in ('name', 'id'):
            assert field in school

        assert school.get('name') == new_name
        school_id = school.get('id')

        assert 'location' in response.headers
        assert response.headers.get('location') == f'{self.SCHOOLS_PATH}/{school_id}'

    def test_add_school_bad_name_type(self):
        new_name = 42
        response = self.add_school(name=new_name)
        assert response.status_code == HTTPStatus.BAD_REQUEST
        check_error_response(response, 'Field "name" must be a string')
