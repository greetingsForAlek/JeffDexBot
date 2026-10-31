# Contributing to JeffDexBot
First of all, thank you for being interested in contributing to the JeffDexBot! All contributions, wether that's bug reports, issues, feature suggestions, documentation improvements or code changes are greatly appreciated.

## Getting Started
Before contributing, make sure you have these:
- Go installed on your system
- Git installed on your system
- A Github Account

## Setting up the project
1. Fork the repository.
2. Clone your fork: `git clone https://github.com/greetingsForAlek/JeffDexBot.git`
3. Navigate to the project directory: `cd JeffDexBot`
4. Install the dependencies: `go mod download`
5. Configure your environment variables as described in the README.

**Important:** Never commit your Discord bot token, database credentials or other secrets.

## Reporting Bugs
If you encounter a bug, please open a Github Issue.
When reporting a bug, try to include:

- A clear description of the problem
- Steps to reproduce the issue
- What you expected to happen
- What actually happened
- Relevant error or message logs

Please avoid including sensitive information such as bot tokens.

## Suggesting Features
Feature suggestions are welcome!
When suggesting a feature, explain:

- What the feature would do
- Why it would be useful
- How you imagine it working

Feel free to open a Github Issue to discuss your idea before implementing it.

## Making Changes
When working on a contribution:

- Create a new branch for your changes
- Keep the changes focused on a single feature or fix where possible
- Follow the existing project structure or coding conventions
- Write clear, descriptive commit messages.
- Test your changes before submitting them.

## Code guidelines
JeffDexBot is written in go. Please follow standard Go conventions.

- Use `gofmt` to format code
- Keep functions focused and readable
- Use meaningful variable and function names.
- Handle errors appropriately.
- Avoid introducing unnecessary dependencies.
- Run ` go build ./...` and `go vet ./...` before submitting changes

## Pull Requests
When you're ready to contribute:

1. Push your branch to your fork
2. Open a pull request against the main JeffDexBot repository.
3. Describe what you've changed and why.
4. Link any related issues
5. Be prepared to discuss or revise your changes if requested.

Please keep Pull Requests focused and avoid combining unrelated changes.

## Community Standards

All contributors are expected to follow the project's Code of Conduct.
Be respectful, constructive, and welcoming when interacting with other contributors.

Thank you for helping to make JeffDexBot better!
