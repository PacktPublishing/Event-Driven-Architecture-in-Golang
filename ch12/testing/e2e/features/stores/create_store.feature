Feature: Creating stores

  As a store owner
  I should be able to create new stores

  Scenario: Creating a store called "Waldorf Books"
    Given a valid store owner
    And no store called "Waldorf Books" exists
    When I create the store called "Waldorf Books"
    Then the store is created

  Scenario: Cannot create stores without a name
    Given a valid store owner
    When I create a store called ""
    Then I receive a "the store name cannot be blank" error
