Feature: Zones

  A zone is a geographical area we are interested in tracking the planes

  Scenario: Zone configuration
    Given The configuration contains one zone
    When The application reads the configuration
    Then The application has one zone configured
