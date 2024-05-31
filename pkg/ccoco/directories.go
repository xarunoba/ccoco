package ccoco

import "errors"

type Directories struct {
	Root, Configs, Preflights string
}

func (d *Directories) CheckState() error {
	if d.Root == "" {
		return errors.New("root directory is empty")
	}
	if d.Configs == "" {
		return errors.New("configs directory is empty")
	}
	if d.Preflights == "" {
		return errors.New("preflights directory is empty")
	}
	return nil
}
