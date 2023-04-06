# Pinned REST API Service

This is a simple REST API service that serves pinned repositories for users and organizations on GitHub. The service is built using the Go and utilizes the GitHub GraphQL API to retrieve pinned repositories for a given user. It useful for showing pinned repositories on your website or blog.

## Usage

Before using this service, you will need a personal access token from GitHub. You can create a new token by following the instructions [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).

You need scopes `public_repo` and `read:org` to use this service.

Once you have a token, set it as an environment variable:

```bash
export GITHUB_TOKEN=your_personal_access_token
```

- `/user/:username` - Returns a JSON response containing the pinned repositories for the specified GitHub user.
- `/org/:orgname` - Returns a JSON response containing the pinned repositories for the specified GitHub organization.

For example

```bash
curl http://localhost:8080/user/joeychilson

curl http://localhost:8080/org/microsoft
```

This will return a JSON response containing the pinned repositories for the specified user or organization.

```json
[
  {
    "name": "vscode",
    "description": "Visual Studio Code",
    "url": "https://github.com/microsoft/vscode",
    "forkCount": 25260,
    "stargazerCount": 144592,
    "language": "TypeScript",
    "updatedAt": "2023-04-06T05:27:36Z",
    "createdAt": "2015-09-03T20:23:38Z"
  }
]
```

## Rate Limiting

The service includes rate limiting functionality to prevent abuse of the GitHub API. If the rate limit has been reached, the service will return a 429 Too Many Requests response with a Retry-After header indicating when the rate limit will reset.

You can view the current rate limit status by checking the headers in the response. The following headers are included:

- `X-RateLimit-Limit` - The maximum number of requests you're permitted to make per hour.
- `X-RateLimit-Remaining` - The number of requests remaining in the current rate limit window.
- `X-RateLimit-Reset` - The time at which the current rate limit window resets in UTC epoch seconds.

GitHub's graphql rate limit is 5000 points per hour. This service uses 1 point per request, so it should be able to handle 1.3 requests per second without hitting the rate limiter.
