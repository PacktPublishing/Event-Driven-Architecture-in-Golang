Feature: Starting baskets

  As a customer I can start new shopping baskets

  Scenario: Create a basket
    Given I am a registered customer
    When I start a new basket
    Then I expect the basket was started
