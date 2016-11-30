# ical2org
An iCalendar to org-mode converter written in Go.

Here is a sample event:

    BEGIN:VEVENT
    UID:uid1@example.com
    DTSTAMP:19970714T170000Z
    ORGANIZER;CN=John Doe:MAILTO:john.doe@example.com
    DTSTART:20161122T170000Z
    DTEND:20161127T035959Z
    SUMMARY:Bastille Day Party
    DESCRIPTION:Hello, World!!!
    END:VEVENT

and the corresponding generated node:

    * Bastille Day Party
      <2016-11-22 Tue 17:00>--<2016-11-27 Sun 03:59>
      :PROPERTIES:
      :ID: uid1@example.com
      :END:
    
      Hello, World!!!

Here are how event properties are mapped onto the Org node:
- The `SUMMARY` is mapped onto the node name;
- the `DESCRIPTION` to the node body;
- the `UID` to the `ID` property;
- the `STATUS` and `LOCATION` to the corresponding properties;
- and the `DTSTART` and `DTEND` are used to bind a date range to the node.
