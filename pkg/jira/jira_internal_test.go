package jira

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildSearchString(t *testing.T) {
	tests := map[string]struct {
		flowType    string
		projectKey  string
		columns     []string
		currentUser bool
		want        string
	}{
		"Kanban without user assigned": {flowType: "kanban", projectKey: "TEST", columns: []string{"ToDo", "InProgress"}, currentUser: false, want: "project = TEST AND status in (\"ToDo\", \"InProgress\") AND type != Epic ORDER BY Rank DESC"},
		"Kanban with user assigned":    {flowType: "kanban", projectKey: "TEST", columns: []string{"ToDo", "InProgress"}, currentUser: true, want: "project = TEST AND status in (\"ToDo\", \"InProgress\") AND type != Epic AND assignee=currentUser() ORDER BY Rank DESC"},
		"Simple without user assigned": {flowType: "simple", projectKey: "TEST", columns: []string{"ToDo", "InProgress"}, currentUser: false, want: "project = TEST AND status in (\"ToDo\", \"InProgress\") AND type != Epic ORDER BY Rank DESC"},
		"Simple with user assigned":    {flowType: "simple", projectKey: "TEST", columns: []string{"ToDo", "InProgress"}, currentUser: true, want: "project = TEST AND status in (\"ToDo\", \"InProgress\") AND type != Epic AND assignee=currentUser() ORDER BY Rank DESC"},
		"Scrum without user assigned":  {flowType: "scrum", projectKey: "TEST", columns: []string{"ToDo", "InProgress"}, currentUser: false, want: "project = TEST AND status in (\"ToDo\", \"InProgress\") AND type != Epic AND sprint in openSprints()"},
		"Scrum with user assigned":     {flowType: "scrum", projectKey: "TEST", columns: []string{"ToDo", "InProgress"}, currentUser: true, want: "project = TEST AND status in (\"ToDo\", \"InProgress\") AND type != Epic AND assignee=currentUser() AND sprint in openSprints()"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := buildSearchString(tc.flowType, tc.projectKey, tc.columns, tc.currentUser)

			if !cmp.Equal(result, tc.want) {
				t.Errorf("Result string '%+v' is different from expected one '%+v'", result, tc.want)
			}
		})
	}
}
