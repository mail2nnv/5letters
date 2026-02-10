package main

import "testing"

func Test_rule_match(t *testing.T) {
	tests := []struct {
		// Named input parameters for receiver constructor.
		arg string
		// Named input parameters for target function.
		w    string
		want bool
	}{
		{"ПИЛКА|20001", "ПОЖАР", true},
		{"ПИЛКА|20001", "ПАПАХ", true},
		{"ПИЛКА|20001", "ОТПАД", false}, // wrong pos `П`
		{"ПИЛКА|20001", "ПОЧТА", false}, // wrong pos `A`
		{"ПИЛКА|20001", "ПАРТА", false}, // wrong pos `A` (last)
		{"ПИЛКА|20001", "ПЕРСИ", false}, // missed `A`
		{"ПИЛКА|20001", "ПАРИЯ", false}, // `И` deprecated
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			r, err := makeRule(tt.arg)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := r.match(tt.w)
			if got != tt.want {
				t.Errorf("match() = %v, want %v", got, tt.want)
			}
		})
	}
}
