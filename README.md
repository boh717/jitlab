## Configuration Prerequisites

Start with the following json file (e.g. ~/.jitlab.json)

```json
{
  "gitlab": {
    "baseurl": <gitlab-url>,
    "groupid": <gitlab-group-id>,
    "token": <gitlab-token>
  },
  "jira": {
    "baseurl": <jira-url>,
    "token": <jira-token>,
    "username": <jira-username>
  },
  "branchPrefix": <your-branch-prefix>,
  "branchSuffix": <your-branch-suffix>,
  "keyCommitSeparator": <your-separator>
}
```

where:
- `<gitlab-token>` is a token with `api` permissions and you can issue one here https://gitlab.com/profile/personal_access_tokens (or similar URL if you are on-premise)
- `<gitlab-group-id>` is the ID of the group your project belongs to
- `<jira-token>` can be issued here: https://id.atlassian.com/manage-profile/security/api-tokens
- `branchPrefix` is what you want to be *prefixed* to every branch you create
- `branchSuffix` is what you want to be *appended* to every branch you create
- `keyCommitSeparator` is what you want to separate the jira key and your commit message

## Board Prerequisites

Jitlab works with both kanban and scrum workflows, but on jira there's a third board type (`simple`) which screws things up.

A `simple` board can be both... And its purpose is to help people setting up boards fast. 

Jitlab treats by default `kanban` and `simple` the same way, but if you're using `scrum` workflow and your team's board is `simple`, then jitlab won't work.

In future releases I may fix this problem or simply ask to change your board type to `scrum`.

## Configuration

Once you fulfilled all the requirements, start using jitlab by configuring it.

Run `jitlab config` and follow the questions you'll be asked. You should run this command only once (or if you change the board).

Your `.jitlab.json` will be updated.

## Project init

Every project you want to use jitlab with should be initialized.

Run `jitlab init` to do it. This will create a local `.repo` file with the project information.
For example:
```json
{
    "id": 12345678,
    "name": "Jitlab",
    "description": "An awesome tool",
    "path": "jitlab"
}
```

## Working on tasks

Jitlab will read issues from jira and will create a local git branch according to the jira task title.

Use `jitlab new` to pick up tasks from your chosen columns.

Branches will follow this naming convention `<your-prefix>TEST-12-your-branch-title<your-suffix>`.

## Pushing changes

Jitlab supports `git commit` and will automatically prefix the message with the jira key. One example could be `TEST-12: awesome message` where `:` is your chosen key commit separator.

Run `jitlab commit -m 'awesome message'`

## Creating merge request

Once you're happy with your changes, you can create the merge request issuing `jitlab mr`.
