package di

import "go.uber.org/fx"

// TODO: 後でdiプロジェクトに移す

var constructors = []any{}

func AppendConstructors(adding []any) error {
	constructors = append(constructors, adding...)
	return nil
}

func Constructors() []any {
	return constructors
}

func AsModule(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"modules"`),
	)
}

func Bind[T any](concrete any) any {
	return fx.Annotate(concrete, fx.As(new(T)))
}
