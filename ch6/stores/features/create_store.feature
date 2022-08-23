Feature: Create Store

  As a store owner
  I should be able to create new stores

  Scenario: Creating a store called "Waldorf Books"
    Given a valid store owner
    And no store called "Waldorf Books" exists
    When I create the store called "Waldorf Books"
    Then a store called "Waldorf Books" exists
