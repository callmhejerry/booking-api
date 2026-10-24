import 'package:flutter/material.dart';
import 'package:booking_web/main.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('booking page shows API settings and booking form', (
    tester,
  ) async {
    tester.view.physicalSize = const Size(1000, 1600);
    tester.view.devicePixelRatio = 1;
    addTearDown(tester.view.resetPhysicalSize);
    addTearDown(tester.view.resetDevicePixelRatio);
    await tester.pumpWidget(const BookingWebApp());

    expect(find.text('Find a time that works.'), findsOneWidget);
    expect(find.text('API base URL'), findsOneWidget);
    expect(find.text('API key'), findsOneWidget);
    expect(find.text('Provider'), findsOneWidget);
    expect(find.text('Available slots'), findsOneWidget);
    expect(find.text('Your details'), findsOneWidget);
    expect(find.text('Confirm booking'), findsOneWidget);
    expect(find.text('Bookings'), findsOneWidget);
    expect(find.byKey(const Key('booking-email-search')), findsOneWidget);
    expect(find.byKey(const Key('list-bookings')), findsOneWidget);
    expect(find.byKey(const Key('booking-id-search')), findsOneWidget);
    expect(find.byKey(const Key('find-booking-by-id')), findsOneWidget);

    expect(find.text('Administrative: replace all slots'), findsOneWidget);
    expect(find.byKey(const Key('reset-slots-json')), findsOneWidget);
  });

  testWidgets('reset requires explicit confirmation', (tester) async {
    tester.view.physicalSize = const Size(1000, 1600);
    tester.view.devicePixelRatio = 1;
    addTearDown(tester.view.resetPhysicalSize);
    addTearDown(tester.view.resetDevicePixelRatio);
    await tester.pumpWidget(const BookingWebApp());
    await tester.enterText(find.byType(TextField).first, 'http://localhost');
    final resetButton = find.byKey(const Key('replace-all-slots'));
    await tester.ensureVisible(resetButton);
    await tester.pumpAndSettle();
    await tester.tap(resetButton);
    await tester.pumpAndSettle();

    expect(find.text('Replace all slots?'), findsOneWidget);
    expect(
      find.textContaining(
        'Existing bookings can also be deleted by the database cascade',
      ),
      findsOneWidget,
    );

    await tester.tap(find.text('Cancel'));
    await tester.pumpAndSettle();

    expect(find.text('Replace all slots?'), findsNothing);
  });
}
