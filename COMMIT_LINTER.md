# Rules
type(scope): message

### Type
- **build**: Changes that affect the build system or external dependencies (e.g., upgrading dependencies, modifying build scripts).

  Use cases:
    - Adding a new dependency in package.json.
    - Updating the version of a dependency.
    - Modifying the build configuration (e.g., Webpack, Maven, Gradle).
- **chore:** Routine tasks that don’t modify source code or add new features or bug fixes (e.g., updating scripts, configurations).

  Use Cases:
    - Updating the .gitignore file.
    - Changing the format of documentation without altering content.
    - Renaming files or folders without changing the code inside.
    - Adding a package for internal tools or developer utilities.
- **ci:** Changes related to Continuous Integration configuration and scripts (e.g., modifying .travis.yml, GitHub Actions).

  Use Cases:
    - Adding a new GitHub Actions workflow.
    - Updating CI/CD configuration files.
    - Fixing a CI script to correct a build issue.
- **docs:** Changes related to documentation only.

  Use Cases:
    - Writing or updating README.md.
    - Adding or updating API documentation.
    - Correcting typos in documentation files.
- **feat:** Introducing a new feature to the codebase.

  Use Cases:
    - Adding a new API endpoint.
    - Implementing a new user interface component.
    - Developing a new module or feature in the application.
- **fix:** A bug fix that addresses an issue in the code.

  Use Cases:
    - Fixing a null pointer exception.
    - Correcting a broken link in the application.
    - Resolving an off-by-one error in a loop.
- **perf:** Code changes that improve performance without altering behavior.

  Use Cases:
    - Optimizing a database query.
    - Reducing the load time of a webpage.
    - Improving the efficiency of an algorithm.
- **refactor:** Code changes that neither fix bugs nor add features but improve the code structure.

  Use Cases:
    - Renaming variables for clarity.
    - Reorganizing code into more modular functions.
    - Removing duplicate code by abstracting common functionality.
- **revert:** A commit that undoes a previous commit.

  Use Cases:
    - Reverting a commit that introduced a bug.
    - Undoing a change that caused a failed build.
    - Rolling back a feature that’s not ready for release.
- **style:** Changes that do not affect the meaning of the code (e.g., formatting, whitespace, semi-colons, linting).

  Use Cases:
    - Correcting code indentation.
    - Reformatting code according to style guidelines.
    - Adding or removing spaces or semi-colons.
- **test:** Adding or modifying tests to ensure the code behaves as expected.

  Use Cases:
    - Writing new unit tests for a module.
    - Updating tests to cover edge cases.
    - Fixing broken tests after a refactor.
