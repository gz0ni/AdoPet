import 'dart:io';

import 'package:flutter/widgets.dart';
import 'package:window_manager/window_manager.dart';

import 'ui_fluent/fluent_app.dart';
import 'ui_macos/macos_app.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await windowManager.ensureInitialized();
  await windowManager.waitUntilReadyToShow(
    const WindowOptions(size: Size(900, 620), title: 'Claude Pet Settings', center: true),
    () async {
      await windowManager.show();
      await windowManager.focus();
    },
  );
  runApp(Platform.isMacOS ? const PetMacosApp() : const PetFluentApp());
}
