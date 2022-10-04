Feature: Creating products

  As a store owner
  I should be able to create new store products

  Background:
    Given the store "Nic & Naks" already exists

  Scenario: Creating a product called "Wizard w/ crystal"
    Given a valid store owner
    When I create the product called "Wizard w/ crystal" with price "9.98"
    Then I expect the product was created
