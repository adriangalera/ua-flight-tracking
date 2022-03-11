Feature: Aircraft

  The flight tracking application will detect when a new plane enters or departs the zone

  Scenario: A new plane appears
    Given A zone is configured
    When New plane appears in the zone
    Then The application notifies

  Scenario: A new plane departs
    Given A zone is configured
    When New plane departs in the zone
    Then The application notifies

