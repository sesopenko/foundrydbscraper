package foundrydata

import (
	"html/template"
	"testing"
)

func Test_foundryTagProcessor(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want template.HTML
	}{
		// TODO: Add test cases.
		{
			name: "Thievery disable hidden mechanism",
			args: args{
				text: "@Check[type:thievery|dc:33|name:Disable Hidden Mechanism|traits:action:disable-a-device]",
			},
			want: `<span class="check">Disable Hidden Mechanism (Thievery DC 33)</span>`,
		},
		{
			name: "Simple Perception with description",
			args: args{
				text: `@Check[type:perception|dc:18]{Something Else}`,
			},
			want: `<span class="check">Something Else (Perception DC 18)</span>`,
		},
		{
			name: "Journal Entry with Entry Page",
			args: args{
				text: `@UUID[JournalEntry.3T1M395V6J75OsEp.JournalEntryPage.8Jl0TWoH4iUJzJxK]{Beginning the Adventure}`,
			},
			want: `<a href="/journal_pages/8Jl0TWoH4iUJzJxK.html">Beginning the Adventure</a>`,
		},
		{
			name: "Simple, no description check",
			args: args{
				text: `@Check[type:fortitude|dc:30]`,
			},
			want: `<span class="check">Fortitude (DC 30)</span>`,
		},
		{
			name: "Actor link with desc",
			args: args{
				text: `@Actor[qXT1SQDtGqMkVl7Q]{Shanrigol Heaps (3)}`,
			},
			want: `<a href="/actors/qXT1SQDtGqMkVl7Q.html">Shanrigol Heaps (3)</a>`,
		},
		//{
		//	name: "Unhandled things",
		//	args: args{
		//		text: `@Macro[m5Crw7ba08oqJdXc]{E06 - Bridge}`,
		//	},
		//	want: "E06 - Bridge",
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := foundryTagProcessor(tt.args.text); got != tt.want {
				t.Errorf("foundryTagProcessor() = %v, want %v", got, tt.want)
			}
		})
	}
}
