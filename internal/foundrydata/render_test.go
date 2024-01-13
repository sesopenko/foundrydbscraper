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
		{
			name: "Spell Effect",
			args: args{
				text: `@UUID[Compendium.pf2e.spell-effects.Item.NXzo2kdgVixIZ2T1]{Spell Effect: Apex Companion}`,
			},
			want: "Spell Effect: Apex Companion",
		},
		{
			name: "Condition",
			args: args{
				text: `@UUID[Compendium.pf2e.conditionitems.Item.dfCMdR4wnpbYNTix]{Stunned 1}`,
			},
			want: `<span class="condition">Stunned 1</span>`,
		},
		{
			name: "Condition without Item",
			args: args{
				text: `@UUID[Compendium.pf2e.conditionitems.MIRkyAjyBeXivMa7]{Enfeebled 2}`,
			},
			want: `<span class="condition">Enfeebled 2</span>`,
		},
		{
			name: "Roll with damage type",
			args: args{
				text: `[[/r (1d10+6)[bleed]]]`,
			},
			want: `<span class="roll">1d10+6 bleed</span>`,
		},
		{
			name: "Item Link",
			args: args{
				text: `@Item[c5oQP02ulBQ7nIVs]{+1 Morningstar}`,
			},
			want: `<a href="/items/c5oQP02ulBQ7nIVs.html">+1 Morningstar</a>`,
		},
		{
			name: "roll with desc",
			args: args{
				text: `[[/r 1d4 #rounds]]{1d4 rounds}`,
			},
			want: `<span class="roll">1d4 rounds</span>`,
		},
		{
			name: "Simple roll",
			args: args{
				text: `[[/r 2d10]]`,
			},
			want: `<span class="roll">2d10</span>`,
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
