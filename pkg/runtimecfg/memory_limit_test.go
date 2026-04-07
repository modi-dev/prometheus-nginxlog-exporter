package runtimecfg

import "testing"

func TestParseCgroupMemoryLimitValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		limit     int64
		hasLimit  bool
		shouldErr bool
	}{
		{
			name:     "valid limit",
			input:    "268435456\n",
			limit:    268435456,
			hasLimit: true,
		},
		{
			name:     "cgroup v2 unlimited",
			input:    "max",
			hasLimit: false,
		},
		{
			name:     "cgroup v1 unlimited sentinel",
			input:    "9223372036854771712",
			hasLimit: false,
		},
		{
			name:     "zero limit",
			input:    "0",
			hasLimit: false,
		},
		{
			name:      "invalid number",
			input:     "not-a-number",
			shouldErr: true,
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			limit, hasLimit, err := parseCgroupMemoryLimitValue(testCase.input)
			if testCase.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if hasLimit != testCase.hasLimit {
				t.Fatalf("unexpected hasLimit; expected %v, got %v", testCase.hasLimit, hasLimit)
			}

			if limit != testCase.limit {
				t.Fatalf("unexpected limit; expected %d, got %d", testCase.limit, limit)
			}
		})
	}
}
