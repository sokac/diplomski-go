describe('Fibonacci test', function() {
    it("Ispravna forma", function() {
        cy.visit('/');

        cy.get('#number').type('11');
        cy.get('form').contains('Računaj').click();
        cy.contains('89');
    });

    it("Neispravna forma", function() {
        cy.visit('/');

        cy.get('#number').type("-1");
        cy.get('form').contains('Računaj').click();
        cy.contains('Parametar n mora biti prirodni broj.');
    });
})

