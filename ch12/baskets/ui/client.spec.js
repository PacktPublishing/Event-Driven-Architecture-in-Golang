const {PactV3, MatchersV3} = require('@pact-foundation/pact');
const chai = require('chai');
const expect = chai.expect;
const axios = require("axios");

const {Client} = require('./client');

describe('Baskets UI', () => {
  let provider;

  before(async () => {
    provider = new PactV3({
      consumer: 'baskets-ui',
      provider: 'baskets-api',
      logLevel: "warn",
      spec: 4,
    });
  });

  context('calling startBasket', () => {
    describe('with a customer ID', () => {
      let basketId;

      before(() => {
        basketId = 'basket-id';

        provider.uponReceiving('a request to start a basket')
          .withRequest({
            method: 'POST',
            path: '/api/baskets',
            body: {
              customerId: 'customer-id'
            },
            headers: {Accept: 'application/json'},
          })
          .willRespondWith({
            body: MatchersV3.like({
              id: MatchersV3.uuid(),
            }),
            headers: {'Content-Type': 'application/json'},
            status: 200,
          });
      });

      it('should start a new basket for the customer', () => {
        return provider.executeTest((mockServer) => {
          const client = new Client(mockServer.url);
          return client.startBasket('customer-id')
            .then((response) => {
              expect(response.data).to.have.property('id');
            });
        });
      });
    });
    describe('without any customer ID', () => {
      before(() => {
        provider.uponReceiving('a request to start a basket without a customerId')
          .withRequest({
            method: 'POST',
            path: '/api/baskets',
            body: {
              customerId: ''
            },
            headers: {Accept: 'application/json'},
          })
          .willRespondWith({
            body: MatchersV3.like({
              "message": "the customer id cannot be blank",
            }),
            headers: {'Content-Type': 'application/json'},
            status: 400,
          });
      });

      it('should not start a new basket for the customer', async () => {
        return provider.executeTest((mockServer) => {
          const client = new Client(mockServer.url);
          return client.startBasket('')
            .catch(({response}) => {
              expect(response.status).to.eq(400);
            });
        });
      });
    });
  });

  context('calling addItem', () => {
    let basketId;
    let productId;

    describe('with a valid product and quantity', () => {
      before(() => {
        basketId = 'basket-id';
        productId = 'product-id';

        provider.given('a store exists')
          .given('a product exists', {id: productId})
          .given('a basket exists', {id: basketId})
          .uponReceiving('a request to add a product')
          .withRequest({
            method: 'PUT',
            path: `/api/baskets/${basketId}/addItem`,
            body: {
              productId: productId,
              quantity: 1,
            },
            headers: {Accept: 'application/json'},
          })
          .willRespondWith({
            body: MatchersV3.equal({}),
            headers: {'Content-Type': 'application/json'},
            status: 200,
          });
      });

      it('should add the item to the basket', () => {
        return provider.executeTest((mockServer) => {
          const client = new Client(mockServer.url);
          return client.addItem(basketId, productId, 1)
            .then((response) => {
              expect(response.status).to.eq(200);
            });
        });
      });
    });
    describe('with a valid product and zero quantity', () => {
      before(() => {
        basketId = 'basket-id';
        productId = 'product-id';

        provider.given('a store exists')
          .given('a product exists', {id: productId})
          .given('a basket exists', {id: basketId})
          .uponReceiving('a request to add a product with a negative quantity')
          .withRequest({
            method: 'PUT',
            path: `/api/baskets/${basketId}/addItem`,
            body: {
              productId: productId,
              quantity: -1,
            },
            headers: {Accept: 'application/json'},
          })
          .willRespondWith({
            body: MatchersV3.like({
              'message': 'the item quantity cannot be negative',
            }),
            headers: {'Content-Type': 'application/json'},
            status: 400,
          });
      });

      it('should return an error message', () => {
        return provider.executeTest((mockServer) => {
          const client = new Client(mockServer.url);
          return client.addItem(basketId, productId, -1)
            .catch(({response}) => {
              expect(response.status).to.eq(400);
            });
        });
      });
    });
    describe('with an unknown product', () => {
      before(() => {
        basketId = 'basket-id';
        productId = 'product-id';

        provider.given('a basket exists', {id: basketId})
          .uponReceiving('a request to add a product')
          .withRequest({
            method: 'PUT',
            path: `/api/baskets/${basketId}/addItem`,
            body: {
              productId: productId,
              quantity: 1,
            },
            headers: {Accept: 'application/json'},
          })
          .willRespondWith({
            body: MatchersV3.like({
              message: `product with id: \`${productId}\` does not exist`
            }),
            headers: {'Content-Type': 'application/json'},
            status: 404,
          });
      });

      it('returns an error message', () => {
        return provider.executeTest((mockServer) => {
          const client = new Client(mockServer.url);
          return client.addItem(basketId, productId, 1)
            .catch(({response}) => {
              expect(response.status).to.eq(404);
            });
        });
      });
    });
  });
})
