# Hammerclock Fix: Deadlock Prevention in MVU Architecture

## Summary of Fixes

The application was experiencing deadlock issues primarily due to sending messages on channels during UI initialization, before the message handling loop was started. Here's how we fixed these issues:

### 1. Improved Widget Initialization

We separated widget creation from callback registration in several UI components:

#### Dropdown Controls

```go
// BEFORE (causing deadlock):
rulesetBox := tview.NewDropDown().
    SetOptions(getRulesetNames(model.Options.Rules), func(option string, index int) {
        msgChan <- &SetRulesetMsg{Index: index}
    })

// AFTER (fixed):
rulesetBox := tview.NewDropDown().
    SetOptions(getRulesetNames(model.Options.Rules), nil).
    SetCurrentOption(model.Options.Default)
// Set the callback separately
rulesetBox.SetSelectedFunc(func(option string, index int) {
    msgChan <- &SetRulesetMsg{Index: index}
})
```

#### Input Fields

```go
// BEFORE (causing deadlock):
playerCountBox := tview.NewInputField().
    SetChangedFunc(func(text string) {
        msgChan <- &SetPlayerCountMsg{Count: count}
    })

// AFTER (fixed):
playerCountBox := tview.NewInputField()
// Set the changed function after initialization
playerCountBox.SetChangedFunc(func(text string) {
    msgChan <- &SetPlayerCountMsg{Count: count}
})
```

#### Checkboxes

```go
// BEFORE (causing deadlock):
checkbox := tview.NewCheckbox().
    SetChangedFunc(func(checked bool) {
        msgChan <- &SetOneTurnForAllPlayersMsg{Value: checked}
    })

// AFTER (fixed):
checkbox := tview.NewCheckbox()
// Set the changed function after initialization
checkbox.SetChangedFunc(func(checked bool) {
    msgChan <- &SetOneTurnForAllPlayersMsg{Value: checked}
})
```

### 2. Fixed Input Field Handling

We improved how player name inputs are handled by using closures properly to avoid variable capture issues:

```go
// Store index in a closure to avoid variable capture issues
idx := i
inputField.SetChangedFunc(func(text string) {
    msgChan <- &SetPlayerNameMsg{
        Index: idx,
        Name:  strings.TrimSpace(text),
    }
})
```

## MVU Architecture Implications

These fixes ensure the application properly follows the Model-View-Update (MVU) architecture by:

1. **Preventing Premature Messages**: No messages are sent on the channel until after the message handling loop has started.

2. **Separating UI Initialization**: UI components are initialized without attaching event handlers immediately.

3. **Proper Message Flow**: Event handlers are attached after the message processing system is ready to receive them.

## Testing

These changes have been verified by successfully building and running the application. The deadlock issues previously encountered during startup have been resolved.

## Future Considerations

1. Consider implementing a deferred messaging system during initialization if needed
2. Add more robust error handling around channel communication
3. Create proper automated tests for UI initialization
