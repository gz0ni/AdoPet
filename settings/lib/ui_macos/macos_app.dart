import 'package:flutter/cupertino.dart';
import 'package:macos_ui/macos_ui.dart';

import '../pages/pages.dart';

class PetMacosApp extends StatefulWidget {
  const PetMacosApp({super.key});

  @override
  State<PetMacosApp> createState() => _PetMacosAppState();
}

class _PetMacosAppState extends State<PetMacosApp> {
  int _index = 0;

  @override
  Widget build(BuildContext context) {
    return MacosApp(
      title: 'Claude Pet Settings',
      debugShowCheckedModeBanner: false,
      home: MacosWindow(
        sidebar: Sidebar(
          minWidth: 200,
          builder: (context, controller) => SidebarItems(
            currentIndex: _index,
            onChanged: (i) => setState(() => _index = i),
            items: [
              for (final p in SettingsPage.values) SidebarItem(label: Text(p.title)),
            ],
          ),
        ),
        child: MacosScaffold(
          toolBar: ToolBar(title: Text(SettingsPage.values[_index].title)),
          children: [ContentArea(builder: (context, _) => const SizedBox.shrink())],
        ),
      ),
    );
  }
}
