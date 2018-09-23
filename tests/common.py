import os
from urllib.parse import urlparse, parse_qs

import requests

def check_error_response(response, expected_msg):
    """Checks an error response for structure and content

    This function uses assert for its checks.

    params:
        response: A requests.Response object
        expected_msg: A string with the message expected in the 'message' property.
    """
    body = response.json()
    assert 'message' in body
    message = body.get('message')
    assert message == expected_msg, f'message in body "{message}" does not match expected "{expected_msg}"'

def check_page_link(links, link_name, offset, limit):
    """Check the validity of a page link

    params:
        links: A dictionary with names as keys and URLs as values
        link_name: The link to check
        offset: The expected offset query parameter in the URL
        limit: The expected limt query paramter in the URL
    """
    assert link_name in links
    link = links.get(link_name)
    parse_result = urlparse(link)
    query = parse_qs(parse_result.query)
    assert 'limit' in query
    q_limit = query.get('limit')[0]
    assert_msg = f'query limit is {q_limit}, expected {limit}'
    assert q_limit == str(limit), assert_msg

    assert 'offset' in query
    q_offset = query.get('offset')[0]
    assert_msg = f'query offset is {q_offset}, expected {offset}'
    assert q_offset == str(offset), assert_msg

class TestSchoolsAPI:
    BASE_URL = os.environ.get('BASE_URL', 'http://localhost:8080')
    SCHOOLS_PATH = f'{BASE_URL}/schools'

    def list_schools(self, offset=0, limit=100):
        """Make a request to the ListSchools operation

        params:
            offset: The offset into the collection
            limit: The maximum number of items in the collection to return per page.

        returns:
            A requests.Response object
        """
        response = requests.get(
            url=self.SCHOOLS_PATH,
            params={
                "offset": offset,
                "limit": limit,
            },
        )
        return response

    def list_schools_custom(self, params=None):
        """Make a request to the ListSchools operation with custom parameters

        params:
            params: Dictionary of parameters to pass directly into the requests.get call

        returns:
            A requests.Response object
            """
        response = requests.get(
            url=self.SCHOOLS_PATH,
            params=params,
        )
        return response
