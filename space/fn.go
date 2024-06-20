package space

import "context"

const AttrKey = "$d_SPACE_ATTR"
const GlobalNameKey = "$d_GLOBAL_SPACE_NAME"

// region core functions

func GetSpace(ctx context.Context) (Space, bool) {
	space, ok := ctx.Value(AttrKey).(Space)
	return space, ok
}

func NewSpace(parent context.Context, name string) context.Context {
	parentSpace, ok := GetSpace(parent)
	if ok {
		return context.WithValue(parent, AttrKey, newSpaceWithParent(name, parentSpace))
	} else {
		return context.WithValue(parent, AttrKey, newSpace(name))
	}
}

func SetInNewSpace(ctx context.Context, name, key string, val any) context.Context {
	return Set(NewSpace(ctx, name), key, val)
}

func Set(ctx context.Context, key string, val any) context.Context {
	space, ok := GetSpace(ctx)
	if ok {
		space.Set(key, val)
	}
	return ctx
}

func Has(ctx context.Context, key string) bool {
	space, ok := GetSpace(ctx)
	if !ok {
		return false
	}
	return space.Has(key)
}

func Get(ctx context.Context, key string) (any, bool) {
	space, ok := GetSpace(ctx)
	if ok {
		return space.Get(key)
	}
	return nil, false
}

// endregion

// region GetXXX helper function

func GetString(ctx context.Context, key string) (string, bool) {
	val, ok := Get(ctx, key)
	if ok {
		return val.(string), true
	}
	return "", false
}

func GetBool(ctx context.Context, key string) (bool, bool) {
	val, ok := Get(ctx, key)
	if ok {
		return val.(bool), true
	}
	return false, false
}

func GetRune(ctx context.Context, key string) (rune, bool) {
	val, ok := Get(ctx, key)
	if ok {
		return val.(rune), true
	}
	return 0, false
}

func GetInt(ctx context.Context, key string) (int, bool) {
	val, ok := Get(ctx, key)
	if ok {
		return val.(int), true
	}
	return 0, false
}

func GetInt64(ctx context.Context, key string) (int64, bool) {
	val, ok := Get(ctx, key)
	if ok {
		return val.(int64), true
	}
	return 0, false
}

func GetFloat32(ctx context.Context, key string) (float32, bool) {
	val, ok := Get(ctx, key)
	if ok {
		return val.(float32), true
	}
	return 0, false
}

func GetFloat64(ctx context.Context, key string) (float64, bool) {
	val, ok := Get(ctx, key)
	if ok {
		return val.(float64), true
	}
	return 0, false
}

func GetComplex64(ctx context.Context, key string) (complex64, bool) {
	val, ok := Get(ctx, key)
	if ok {
		return val.(complex64), true
	}
	return 0, false
}

func GetComplex128(ctx context.Context, key string) (complex128, bool) {
	val, ok := Get(ctx, key)
	if ok {
		return val.(complex128), true
	}
	return 0, false
}

// endregion
