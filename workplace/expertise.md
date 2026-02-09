# Code Quality Expertise

## Core Principles

### SOLID Principles
- **S**ingle Responsibility: Each class/function does one thing well
- **O**pen/Closed: Open for extension, closed for modification
- **L**iskov Substitution: Subtypes must be substitutable for base types
- **I**nterface Segregation: Many specific interfaces > one general interface
- **D**ependency Inversion: Depend on abstractions, not concretions

### Clean Code Fundamentals
- **Meaningful Names**: Variables, functions, classes should reveal intent
- **Small Functions**: 5-20 lines, one level of abstraction, one purpose
- **Clear Abstractions**: Each layer should make sense independently
- **No Duplication**: DRY (Don't Repeat Yourself) religiously
- **Simple Design**: KISS (Keep It Simple, Stupid) always

## Code Quality Standards

### Naming Conventions
- GOOD: calculateMonthlyPayment(), isUserAuthenticated(), getUserById()
- BAD: calc(), flag, doStuff(), temp, data
- Variables: Nouns (clear, descriptive)
- Functions: Verbs (action-oriented)
- Classes: Nouns (singular, specific)
- Constants: SCREAMING_SNAKE_CASE
- Booleans: is/has/can prefix

### Function Design
- Max 3-4 parameters (use objects for more)
- Single responsibility per function
- No side effects (function name should reflect all it does)
- Return early to avoid deep nesting
- Use guard clauses for validation
- Avoid flag parameters (split into separate functions)

### Code Organization
- Group related code together
- Order methods: public → private, high-level → low-level
- Keep files under 300 lines
- Keep classes focused (< 10 methods ideally)
- Use modules/namespaces to organize

### Comments & Documentation
- **Don't comment WHAT**, comment **WHY**
- GOOD: // Using exponential backoff to avoid overwhelming the API
- BAD: // Loop through array
- Update comments when updating code
- Remove commented-out code (use version control)
- Add docstrings for public APIs

### Error Handling
- Never ignore errors
- Fail fast and loudly
- Provide helpful error messages
- Use exceptions for exceptional cases
- Return error codes for expected failures
- Log errors with context

### Testing
- Write tests first (TDD when possible)
- Test behavior, not implementation
- One assertion per test (when practical)
- Use descriptive test names: test_calculateTax_withValidInput_returnsCorrectAmount()
- Cover edge cases: null, empty, negative, boundary values
- Keep tests independent and fast

## Code Review Checklist

Before submitting code:
- [ ] All functions under 20 lines?
- [ ] All names self-explanatory?
- [ ] No magic numbers (use named constants)?
- [ ] No code duplication?
- [ ] Error cases handled?
- [ ] Unit tests written and passing?
- [ ] Complex logic commented (why, not what)?
- [ ] Could a junior developer understand this?

## Refactoring Patterns

Common improvements to suggest:
1. **Extract Method**: Long function → multiple small functions
2. **Extract Variable**: Complex expression → named variable
3. **Rename**: Unclear name → descriptive name
4. **Remove Duplication**: Repeated code → shared function
5. **Simplify Conditionals**: Complex if → early returns or extract method
6. **Replace Magic Numbers**: Hardcoded values → named constants
7. **Introduce Parameter Object**: Many parameters → single object

## Language-Specific Best Practices

Always follow idioms for the language in use:
- **Python**: PEP 8, list comprehensions, context managers
- **JavaScript**: Modern ES6+, const/let not var, arrow functions
- **Go**: Error handling, goroutines, defer
- **Java**: Streams, Optional, builder pattern
- **TypeScript**: Type safety, interfaces, generics
- **Rust**: Ownership, Result types, pattern matching

## Metrics to Watch
- **Cyclomatic Complexity**: Keep under 10
- **Code Coverage**: Aim for 80%+
- **Code Duplication**: Under 3%
- **Function Length**: Under 20 lines
- **File Length**: Under 300 lines
- **Class Size**: Under 10 methods
