package redl

import (
	"testing"

	"github.com/go-test/deep"
)

func TestExtractTags(t *testing.T) {
	testCases := []struct {
		desc        string
		input       string
		outputTags  []string
		outputTitle string
	}{
		{
			desc:        "Empty String",
			input:       "",
			outputTags:  []string{},
			outputTitle: "",
		},
		{
			desc:        "No Tags",
			input:       "No Tags",
			outputTags:  []string{},
			outputTitle: "No Tags",
		},
		{
			desc:        "Single at Beginning",
			input:       "[Tag1]Single at Beginning",
			outputTags:  []string{"Tag1"},
			outputTitle: "Single at Beginning",
		},
		{
			desc:        "Single at End",
			input:       "Single at End[Tag1]",
			outputTags:  []string{"Tag1"},
			outputTitle: "Single at End",
		},
		{
			desc:        "Two at Beginning",
			input:       "[Tag1][Tag2]Two at Beginning",
			outputTags:  []string{"Tag1", "Tag2"},
			outputTitle: "Two at Beginning",
		},
		{
			desc:        "Two at End",
			input:       "Two at End[Tag2][Tag1]",
			outputTags:  []string{"Tag1", "Tag2"},
			outputTitle: "Two at End",
		},
		{
			desc:        "Both Ends",
			input:       "[Tag1][Tag2]Both Ends[Tag2][Tag1]",
			outputTags:  []string{"Tag1", "Tag2", "Tag1", "Tag2"},
			outputTitle: "Both Ends",
		},
		{
			desc:        "Ignore Empty",
			input:       "[]Ignore Empty[]",
			outputTags:  []string{},
			outputTitle: "Ignore Empty",
		},
		{
			desc:        "Trim Spaces",
			input:       "   [ Tag 1] Trim Spaces [ Tag 2]   ",
			outputTags:  []string{"Tag 1", "Tag 2"},
			outputTitle: "Trim Spaces",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actualTags, actualTitle := ExtractTags(tC.input)
			if diff := deep.Equal(tC.outputTags, actualTags); diff != nil {
				t.Error(diff)
			}
			if diff := deep.Equal(tC.outputTitle, actualTitle); diff != nil {
				t.Error(diff)
			}
		})
	}
}
