workflow "New workflow" {
  on = "push"
  resolves = ["go"]
}

action "go" {
  uses = "go"
  runs = "run"
  args = "main.go"
}
