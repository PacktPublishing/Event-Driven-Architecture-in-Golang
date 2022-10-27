Feature: Adding items

  As a customer with a basket I can add items

  Background:
    Given I am a registered customer
    And I start a new basket

  Scenario: Adding items
    Given a store has the following items
      | Name              | Price |
      | Wizard w/ crystal | 9.99  |
    When I add the items
      | Name              | Quantity |
      | Wizard w/ crystal | 10       |
    Then the items are added

  Scenario: Cannot add items that do not exist
    When I add the items
      | Name              | Quantity |
      | Wizard w/ crystal | 10       |
    Then I receive a "product fallback failed: product was not located" error
