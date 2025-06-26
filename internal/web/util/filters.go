package util

import (
	"strings"

	"github.com/quinn/gmail-sorter/pkg/gmailapi"
	"google.golang.org/api/gmail/v1"
)

func DescribeFilterCriteria(criteria *gmail.FilterCriteria) string {
	var result []string

	if criteria.From != "" {
		result = append(result, "from: "+criteria.From)
	}

	if criteria.Subject != "" {
		result = append(result, "subject: "+criteria.Subject)
	}

	if criteria.To != "" {
		result = append(result, "to: "+criteria.To)
	}

	if criteria.HasAttachment {
		result = append(result, "has attachment")
	}

	return strings.Join(result, ", ")
}

func DescribeFilterAction(api *gmailapi.GmailAPI, action *gmail.FilterAction) string {
	var result []string

	if len(action.AddLabelIds) > 0 {
		for _, id := range action.AddLabelIds {
			result = append(result, "add label: "+GetLabel(api, id))
		}
	}

	if len(action.RemoveLabelIds) > 0 {
		for _, id := range action.RemoveLabelIds {
			result = append(result, "remove label: "+GetLabel(api, id))
		}
	}

	if action.Forward != "" {
		result = append(result, "forward to: "+action.Forward)
	}

	return strings.Join(result, ", ")
}

func GetLabel(api *gmailapi.GmailAPI, id string) string {
	label, err := api.Label(id)

	if err != nil {
		return "N/A"
	}

	if label.Name == "" {
		return "???"
	}

	return label.Name
}
