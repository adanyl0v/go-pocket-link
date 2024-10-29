package redis

import "fmt"

func errConnecting(err error) error {
	return fmt.Errorf("%w (connecting to redis)", err)
}

func errClosing(err error) error {
	return fmt.Errorf("%w (closing redis client)", err)
}

func errKeyDoesNotExist(key string) error {
	return fmt.Errorf("%s does not exist", key)
}
