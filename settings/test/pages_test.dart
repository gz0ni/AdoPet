import 'package:claude_pet_settings/pages/pages.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('every page has a title', () {
    for (final p in SettingsPage.values) {
      expect(p.title, isNotEmpty);
    }
  });
}
