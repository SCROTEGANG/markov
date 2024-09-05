package markov

import (
	"context"
	"testing"
)

var data = []string{
	"Lorem ipsum dolor sit amet",
	"Lorem ipsum duplicate start words",
	"Consectetur adipiscing elit",
	"Quisque tempor, erat vel lacinia imperdiet",
	"Justo nisi fringilla dui",
	"Egestas bibendum eros nisi ut lacus",
	"fringilla dui avait annoncé une rupture avec le erat vel: il n'en est rien…",
	"Fusce tincidunt tempor, erat vel lacinia vel ex pharetra pretium lacinia imperdiet",
}

func TestMarkov(t *testing.T) {
	// TODO: better tests

	ctx := context.Background()
	m := New(
		WithStateSize(2),
	)

	m.AddData(ctx, data)

	res, err := m.Generate(
		ctx,
		WithMaxTries(10),
	)
	if err != nil {
		t.Fatalf("Error generating sentence: %s", err.Error())
	}

	t.Logf("sentence: %s", res.String)
}
