package shell

func (c *ShellCommand) configure() {

}

func (c *ShellCommand) Exit() error {
	if c.Cmd.Process != nil && !c.Cmd.ProcessState.Exited() {
		return c.Cmd.Process.Kill()
	}

	return nil
}
