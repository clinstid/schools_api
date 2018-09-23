from http import HTTPStatus

from common import (
    TestSchoolsAPI,
    check_error_response,
    check_page_link,
)


class TestListSchools(TestSchoolsAPI):
    def test_list_schools_simple(self):
        response = self.list_schools()
        assert response.status_code == HTTPStatus.OK

        schools_collection = response.json()
        for field in ('schools', 'meta', 'links'):
            assert field in schools_collection

        meta = schools_collection.get('meta')
        assert meta is not None
        assert 'total' in meta

        links = schools_collection.get('links')
        assert links is not None
        for field in ('next', 'first', 'last'):
            assert field in links

    def test_list_schools_bad_limits(self):
        bad_limit_msg = 'limit query parameter must be at least 1 and no greater than 100'
        response = self.list_schools(limit=0)
        assert response.status_code == HTTPStatus.BAD_REQUEST
        check_error_response(response, bad_limit_msg)

        response = self.list_schools(limit=-1)
        assert response.status_code == HTTPStatus.BAD_REQUEST
        check_error_response(response, bad_limit_msg)

        response = self.list_schools(limit=1000)
        assert response.status_code == HTTPStatus.BAD_REQUEST
        check_error_response(response, bad_limit_msg)

    def test_list_schools_bad_offset(self):
        response = self.list_schools(offset=-1)
        assert response.status_code == HTTPStatus.BAD_REQUEST
        check_error_response(response, 'offset query parameter must be at least 0')

    def test_list_schools_pagination(self):
        offset = 10
        limit = 10
        response = self.list_schools(offset=offset, limit=limit)
        assert response.status_code == HTTPStatus.OK

        schools_collection = response.json()
        assert 'meta' in schools_collection
        meta = schools_collection.get('meta')
        assert 'total' in meta
        total = meta.get('total')

        assert 'links' in schools_collection
        links = schools_collection.get('links')
        check_page_link(links, 'next', offset+limit, limit)
        check_page_link(links, 'prev', offset-limit, limit)
        check_page_link(links, 'first', 0, limit)
        check_page_link(links, 'last', (total // limit) * limit , limit)
