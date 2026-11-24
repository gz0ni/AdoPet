import 'package:fluent_ui/fluent_ui.dart';

import '../pages/pages.dart';

class PetFluentApp extends StatefulWidget {
  const PetFluentApp({super.key});

  @override
  State<PetFluentApp> createState() => _PetFluentAppState();
}

class _PetFluentAppState extends State<PetFluentApp> {
  int _index = 0;

  @override
  Widget build(BuildContext context) {
    return FluentApp(
      title: 'Claude Pet Settings',
      debugShowCheckedModeBanner: false,
      home: NavigationView(
        pane: NavigationPane(
          selected: _index,
          onChanged: (i) => setState(() => _index = i),
          items: [
            for (final p in SettingsPage.values)
              PaneItem(
                icon: const Icon(FluentIcons.settings),
                title: Text(p.title),
                body: ScaffoldPage(header: PageHeader(title: Text(p.title))),
              ),
          ],
        ),
      ),
    );
  }
}
