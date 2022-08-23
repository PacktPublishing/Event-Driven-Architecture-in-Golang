Feature: Register Customer

  Scenario: Registering a new customer
    Given no customer named "John Smith" exists
    When I register a new customer as "John Smith"
    Then I expect the request to succeed
    And expect a customer named "John Smith" to exist
