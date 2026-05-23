#ifdef _WIN32
#include <Windows.h>
#include "hotkey.h"

// Declaración de la función Go (se define en Go via //export)
extern void handleHotkey(int id);

// IDs: 1 toggle, 2 configurable quit, 3 Ctrl+Q quit, 4 Ctrl+, preferences, 5 F1 about
#define HOTKEY_ID_TOGGLE 1
#define HOTKEY_ID_QUIT 2
#define HOTKEY_ID_QUIT_CTRL 3
#define HOTKEY_ID_PREFS 4
#define HOTKEY_ID_ABOUT 5

#define MOD_CONTROL 0x0002
#define VK_Q 0x51
#define VK_OEM_COMMA 0xBC
#define VK_F1 0x70

void setupHotkeys(unsigned int toggleModifier, unsigned int toggleKey, unsigned int exitModifier, unsigned int exitKey) {
    RegisterHotKey(NULL, HOTKEY_ID_TOGGLE, toggleModifier, toggleKey);
    RegisterHotKey(NULL, HOTKEY_ID_QUIT, exitModifier, exitKey);
    RegisterHotKey(NULL, HOTKEY_ID_QUIT_CTRL, MOD_CONTROL, VK_Q);
    RegisterHotKey(NULL, HOTKEY_ID_PREFS, MOD_CONTROL, VK_OEM_COMMA);
    RegisterHotKey(NULL, HOTKEY_ID_ABOUT, 0, VK_F1);

    MSG msg = {0};
    while (GetMessage(&msg, NULL, 0, 0) != 0) {
        if (msg.message == WM_HOTKEY) {
            handleHotkey((int)msg.wParam);
            if (msg.wParam == HOTKEY_ID_QUIT || msg.wParam == HOTKEY_ID_QUIT_CTRL) {
                break;
            }
        }
    }

    UnregisterHotKey(NULL, HOTKEY_ID_TOGGLE);
    UnregisterHotKey(NULL, HOTKEY_ID_QUIT);
    UnregisterHotKey(NULL, HOTKEY_ID_QUIT_CTRL);
    UnregisterHotKey(NULL, HOTKEY_ID_PREFS);
    UnregisterHotKey(NULL, HOTKEY_ID_ABOUT);
}

#endif
