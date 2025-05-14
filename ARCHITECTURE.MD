# Model-View-Update (MVU) Architecture in Hammerclock

Hammerclock is built using the Model-View-Update (MVU) architectural pattern, a unidirectional data flow pattern popularized by the [Elm programming language](https://guide.elm-lang.org/architecture/).

## Overview

The MVU pattern consists of three key components:

1. **Model**: Represents the application state
2. **View**: Renders the UI based on the model
3. **Update**: Handles events and updates the model

## Implementation in Hammerclock

### Model

The model represents our entire application state. All data needed by the application is stored in the model. The model is defined in `internal/app/model.go`.

```go
type Model struct {
    Players             []*Player
    Phases              []string
    GameStatus          GameStatus
    CurrentScreen       string
    GameStarted         bool
    Options             Options
    CurrentColorPalette ColorPalette
}
```

### View

The view is responsible for rendering the UI based on the current model. The view is completely derived from the model. No UI state is stored in the view itself; all state comes from the model. The view is defined in `internal/app/view.go`.

```go
type View struct {
    App                   *tview.Application
    MainFlex              *tview.Flex
    PlayerPanelsContainer *tview.Flex
    PlayerPanels          []*tview.Flex
    TopMenu               *tview.TextView
    BottomMenu            *tview.TextView
    StatusPanel           *tview.Flex
    ClockDisplay          *tview.TextView
    OptionsScreen         *tview.Grid
    AboutScreen           *tview.Flex
    MessageChan           chan<- Message
}
```

The `Render` method updates the UI based on the current model:

```go
func (v *View) Render(model *Model) {
    // Update UI components based on the model
}
```

### Update

The update component handles all events and state changes. It receives messages representing user actions or system events and returns a new model. This ensures that all state changes go through a single pipeline. The update logic is defined in `internal/app/update.go`.

```go
func Update(msg Message, model Model) (Model, Command) {
    // Process the message and return an updated model
}
```

## Unidirectional Data Flow

The data in Hammerclock flows in one direction:

1. The **Model** stores the state
2. The **View** renders the UI based on the model
3. User actions trigger **Messages**
4. Messages are processed by the **Update** function
5. The Update function returns a new **Model**
6. The cycle repeats with the new model

This unidirectional flow makes the application more predictable, testable, and easier to reason about.

## Message Passing

Communication between components is done through a message system. When a user interaction occurs, a message is created and sent to the update function. The update function processes the message and returns a new model. This keeps the components decoupled and ensures that all state changes go through a single pipeline.

```go
// Example message types
type StartGameMsg struct{}
type SwitchTurnsMsg struct{}
type SetRulesetMsg struct{
    Index int
}
```

## Immutable Updates

State changes in Hammerclock are immutable. Rather than modifying the existing model, each update function creates a new copy of the model with the changes applied. This ensures that no side effects occur during updates and makes the application more predictable.

```go
// Example of an immutable update
func handleSetRuleset(msg *SetRulesetMsg, model Model) (Model, Command) {
    newModel := model  // Create a copy
    newModel.Options.Default = msg.Index  // Modify the copy
    newModel.Phases = model.Options.Rules[msg.Index].Phases
    return newModel, NoCommand  // Return the new model
}
```

## Benefits

- **Predictability**: State changes are centralized and follow a clear path
- **Testability**: Functions are pure and can be tested in isolation
- **Maintainability**: Components are decoupled and have clear responsibilities
- **Reasoning**: The application flow is easier to understand
- **Debugging**: State changes are tracked through a single pipeline

## Future Improvements

- Introduce more granular message types for specific state changes
- Implement deeper immutability for nested state objects
- Add middleware for cross-cutting concerns like logging and persistence
- Create a message dispatcher for more complex workflows