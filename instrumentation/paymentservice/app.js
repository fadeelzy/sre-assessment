
require('./otel');

const express = require('express');
const { startPaymentSpan } = require('./otel');

const app = express();
app.use(express.json());

app.post('/pay', (req, res) => {
  const { amount, userId } = req.body || {};
  const span = startPaymentSpan(amount || 0, userId || 'unknown');
  // simulate work
  setTimeout(() => {
    span.end();
    res.json({ status: 'paid' });
  }, 20);
});

const port = process.env.PORT || 3000;
app.listen(port, () => console.log(`paymentservice listening ${port}`));