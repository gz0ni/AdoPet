/// Shared page list; platform layers render the same pages with different widgets.
enum SettingsPage { general, animations, events, behavior, hotkeys, about }

extension SettingsPageTitle on SettingsPage {
  String get title => switch (this) {
        SettingsPage.general => 'General',
        SettingsPage.animations => 'Animations',
        SettingsPage.events => 'Claude events',
        SettingsPage.behavior => 'Behavior',
        SettingsPage.hotkeys => 'Hotkeys',
        SettingsPage.about => 'About',
      };
}
