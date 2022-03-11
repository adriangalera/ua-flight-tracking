Feature: Notification

  The flight tracking application will send notification to console

  Scenario: Console notification for new plane
    Given Console notification is configured
    When New plane appear in the zone
    Then The application notifies the console

  Scenario: Console notification for plane missing
    Given Console notification is configured
    When New plane disappear from the zone
    Then The application notifies the console

  Scenario: Telegram notification for new plane
    Given Telegram notification is configured
    When New plane appear in the zone
    Then The application notifies the console

  Scenario: Telegram notification for plane missing
    Given Telegram notification is configured
    When New plane disappear from the zone
    Then The application notifies the console