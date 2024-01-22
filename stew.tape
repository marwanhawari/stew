Output assets/demo.gif

Require stew

Set Shell "bash"
Set FontSize 24
Set Width 1800
Set Height 600
Set Framerate 24
Set PlaybackSpeed 0.6
Set FontFamily "Menlo"

# Install
Type "stew install sharkdp/fd"
Sleep 1s
Enter@3s
Ctrl+L

# Search
Type "stew search fzf"
Sleep 1s
Enter@2s
Down@300ms 5
Sleep 750ms
Up@300ms 5
Enter@2s
Down@300ms
Enter@1500ms
Type "darwin"
Down@300ms
Enter@3s
Ctrl+L

# Upgrade
Type "stew upgrade fzf"
Sleep 1s
Enter@3s
Ctrl+L

# List
Type "stew list --tags"
Sleep 1s
Enter@3s
Ctrl+L
