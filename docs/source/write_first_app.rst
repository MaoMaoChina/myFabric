<<<<<<< HEAD
Writing Your First Application
==============================

.. note:: If you're not yet familiar with the fundamental architecture of a
          Fabric network, you may want to visit the :doc:`key_concepts` section
          prior to continuing.

          It is also worth noting that this tutorial serves as an introduction
          to Fabric applications and uses simple smart contracts and
          applications. For a more in-depth look at Fabric applications and
          smart contracts, check out our
          :doc:`developapps/developing_applications` section or the
          :doc:`tutorial/commercial_paper`.

In this tutorial we'll be looking at a handful of sample programs to see how
Fabric apps work. These applications and the smart contracts they use are
collectively known as ``FabCar``. They provide a great starting point to
understand a Hyperledger Fabric blockchain. You'll learn how to write an
application and smart contract to query and update a ledger, and how to use a
Certificate Authority to generate the X.509 certificates used by applications
which interact with a permissioned blockchain.

We will use the application SDK --- described in detail in the
:doc:`/developapps/application` topic -- to invoke a smart contract which
queries and updates the ledger using the smart contract SDK --- described in
detail in section :doc:`/developapps/smartcontract`.

We’ll go through three principle steps:

  **1. Setting up a development environment.** Our application needs a network
  to interact with, so we'll get a basic network our smart contracts and
  application will use.

  .. image:: images/AppConceptsOverview.png

  **2. Learning about a sample smart contract, FabCar.**
  We’ll inspect the smart contract to learn about the transactions within them,
  and how they are used by applications to query and update the ledger.

  **3. Develop a sample application which uses FabCar.** Our application will
  use the FabCar smart contract to query and update car assets on the ledger.
  We'll get into the code of the apps and the transactions they create,
  including querying a car, querying a range of cars, and creating a new car.

After completing this tutorial you should have a basic understanding of how an
application is programmed in conjunction with a smart contract to interact with
the ledger hosted and replicated on the peers in a Fabric network.

.. note:: These applications are also compatible with :doc:`discovery-overview`
          and :doc:`private-data/private-data`, though we won't explicitly show
          how to use our apps to leverage those features.

Set up the blockchain network
-----------------------------

.. note:: This next section requires you to be in the ``first-network``
          subdirectory within your local clone of the ``fabric-samples`` repo.

If you've already run through :doc:`build_network`, you will have downloaded
``fabric-samples`` and have a network up and running. Before you run this
tutorial, you must stop this network:

.. code:: bash

  ./byfn.sh down

If you have run through this tutorial before, use the following commands to
kill any stale or active containers. Note, this will take down **all** of your
containers whether they're Fabric related or not.

.. code:: bash

  docker rm -f $(docker ps -aq)
  docker rmi -f $(docker images | grep fabcar | awk '{print $3}')

If you don't have a development environment and the accompanying artifacts for
the network and applications, visit the :doc:`prereqs` page and ensure you have
the necessary dependencies installed on your machine.

Next, if you haven't done so already, visit the :doc:`install` page and follow
the provided instructions. Return to this tutorial once you have cloned the
``fabric-samples`` repository, and downloaded the latest stable Fabric images
and available utilities.

If you are using Mac OS and running Mojave, you will need to `install Xcode
<./tutorial/installxcode.html>`_.

Launch the network
^^^^^^^^^^^^^^^^^^

.. note:: This next section requires you to be in the ``fabcar``
          subdirectory within your local clone of the ``fabric-samples`` repo.

          This tutorial demonstrates the JavaScript versions of the ``FabCar``
          smart contract and application, but the ``fabric-samples`` repo also
          contains Go, Java and TypeScript versions of this sample. To try the
          Go, Java or TypeScript versions, change the ``javascript`` argument
          for ``./startFabric.sh`` below to either ``go``, ``java`` or ``typescript``
          and follow the instructions written to the terminal.

Launch your network using the ``startFabric.sh`` shell script. This command will
spin up a blockchain network comprising peers, orderers, certificate
authorities and more.  It will also install and instantiate a JavaScript version
of the ``FabCar`` smart contract which will be used by our application to access
the ledger. We'll learn more about these components as we go through the
tutorial.

.. code:: bash

  ./startFabric.sh javascript

Alright, you’ve now got a sample network up and running, and the ``FabCar``
smart contract installed and instantiated. Let’s install our application
pre-requisites so that we can try it out, and see how everything works together.

Install the application
^^^^^^^^^^^^^^^^^^^^^^^

.. note:: The following instructions require you to be in the
          ``fabcar/javascript`` subdirectory within your local clone of the
          ``fabric-samples`` repo.

Run the following command to install the Fabric dependencies for the
applications. It will take about a minute to complete:

.. code:: bash

  npm install

This process is installing the key application dependencies defined in
``package.json``. The most important of which is the ``fabric-network`` class;
it enables an application to use identities, wallets, and gateways to connect to
channels, submit transactions, and wait for notifications. This tutorial also
uses the ``fabric-ca-client`` class to enroll users with their respective
certificate authorities, generating a valid identity which is then used by
``fabric-network`` class methods.

Once ``npm install`` completes, everything is in place to run the application.
For this tutorial, you'll primarily be using the application JavaScript files in
the ``fabcar/javascript`` directory. Let's take a look at what's inside:

.. code:: bash

  ls

You should see the following:

.. code:: bash

  enrollAdmin.js  node_modules       package.json  registerUser.js
  invoke.js       package-lock.json  query.js      wallet

There are files for other program languages, for example in the
``fabcar/typescript`` directory. You can read these once you've used the
JavaScript example -- the principles are the same.

If you are using Mac OS and running Mojave, you will need to `install Xcode
<./tutorial/installxcode.html>`_.

Enrolling the admin user
------------------------

.. note:: The following two sections involve communication with the Certificate
          Authority. You may find it useful to stream the CA logs when running
          the upcoming programs by opening a new terminal shell and running
          ``docker logs -f ca.example.com``.

When we created the network, an admin user --- literally called ``admin`` ---
was created as the **registrar** for the certificate authority (CA). Our first
step is to generate the private key, public key, and X.509 certificate for
``admin`` using the ``enroll.js`` program. This process uses a **Certificate
Signing Request** (CSR) --- the private and public key are first generated
locally and the public key is then sent to the CA which returns an encoded
certificate for use by the application. These three credentials are then stored
in the wallet, allowing us to act as an administrator for the CA.

We will subsequently register and enroll a new application user which will be
used by our application to interact with the blockchain.

Let's enroll user ``admin``:

.. code:: bash

  node enrollAdmin.js

This command has stored the CA administrator's credentials in the ``wallet``
directory.

Register and enroll ``user1``
-----------------------------

Now that we have the administrator's credentials in a wallet, we can enroll a
new user --- ``user1`` --- which will be used to query and update the ledger:

.. code:: bash

  node registerUser.js

Similar to the admin enrollment, this program uses a CSR to enroll ``user1`` and
store its credentials alongside those of ``admin`` in the wallet. We now have
identities for two separate users --- ``admin`` and ``user1`` --- and these are
used by our application.

Time to interact with the ledger...

Querying the ledger
-------------------

Each peer in a blockchain network hosts a copy of the ledger, and an application
program can query the ledger by invoking a smart contract which queries the most
recent value of the ledger and returns it to the application.

Here is a simplified representation of how a query works:

.. image:: tutorial/write_first_app.diagram.1.png

Applications read data from the `ledger <./ledger/ledger.html>`_ using a query.
The most common queries involve the current values of data in the ledger -- its
`world state <./ledger/ledger.html#world-state>`_. The world state is
represented as a set of key-value pairs, and applications can query data for a
single key or multiple keys. Moreover, the ledger world state can be configured
to use a database like CouchDB which supports complex queries when key-values
are modeled as JSON data. This can be very helpful when looking for all assets
that match certain keywords with particular values; all cars with a particular
owner, for example.

First, let's run our ``query.js`` program to return a listing of all the cars on
the ledger. This program uses our second identity -- ``user1`` -- to access the
ledger:

.. code:: bash

  node query.js

The output should look like this:

.. code:: json

  Wallet path: ...fabric-samples/fabcar/javascript/wallet
  Transaction has been evaluated, result is:
  [{"Key":"CAR0", "Record":{"colour":"blue","make":"Toyota","model":"Prius","owner":"Tomoko"}},
  {"Key":"CAR1", "Record":{"colour":"red","make":"Ford","model":"Mustang","owner":"Brad"}},
  {"Key":"CAR2", "Record":{"colour":"green","make":"Hyundai","model":"Tucson","owner":"Jin Soo"}},
  {"Key":"CAR3", "Record":{"colour":"yellow","make":"Volkswagen","model":"Passat","owner":"Max"}},
  {"Key":"CAR4", "Record":{"colour":"black","make":"Tesla","model":"S","owner":"Adriana"}},
  {"Key":"CAR5", "Record":{"colour":"purple","make":"Peugeot","model":"205","owner":"Michel"}},
  {"Key":"CAR6", "Record":{"colour":"white","make":"Chery","model":"S22L","owner":"Aarav"}},
  {"Key":"CAR7", "Record":{"colour":"violet","make":"Fiat","model":"Punto","owner":"Pari"}},
  {"Key":"CAR8", "Record":{"colour":"indigo","make":"Tata","model":"Nano","owner":"Valeria"}},
  {"Key":"CAR9", "Record":{"colour":"brown","make":"Holden","model":"Barina","owner":"Shotaro"}}]

Let's take a closer look at this program. Use an editor (e.g. atom or visual
studio) and open ``query.js``.

The application starts by bringing in scope two key classes from the
``fabric-network`` module; ``FileSystemWallet`` and ``Gateway``. These classes
will be used to locate the ``user1`` identity in the wallet, and use it to
connect to the network:

.. code:: bash

  const { FileSystemWallet, Gateway } = require('fabric-network');

The application connects to the network using a gateway:

.. code:: bash

  const gateway = new Gateway();
  await gateway.connect(ccp, { wallet, identity: 'user1' });

This code creates a new gateway and then uses it to connect the application to
the network. ``ccp`` describes the network that the gateway will access with the
identity ``user1`` from ``wallet``. See how the ``ccp`` has been loaded from
``../../basic-network/connection.json`` and parsed as a JSON file:

.. code:: bash

  const ccpPath = path.resolve(__dirname, '..', '..', 'basic-network', 'connection.json');
  const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
  const ccp = JSON.parse(ccpJSON);

If you'd like to understand more about the structure of a connection profile,
and how it defines the network, check out
`the connection profile topic <./developapps/connectionprofile.html>`_.

A network can be divided into multiple channels, and the next important line of
code connects the application to a particular channel within the network,
``mychannel``:

.. code:: bash

  const network = await gateway.getNetwork('mychannel');

Within this channel, we can access the smart contract ``fabcar`` to interact
with the ledger:

.. code:: bash

  const contract = network.getContract('fabcar');

Within ``fabcar`` there are many different **transactions**, and our application
initially uses the ``queryAllCars`` transaction to access the ledger world state
data:

.. code:: bash

  const result = await contract.evaluateTransaction('queryAllCars');

The ``evaluateTransaction`` method represents one of the simplest interaction
with a smart contract in blockchain network. It simply picks a peer defined in
the connection profile and sends the request to it, where it is evaluated. The
smart contract queries all the cars on the peer's copy of the ledger and returns
the result to the application. This interaction does not result in an update the
ledger.

The FabCar smart contract
-------------------------

Let's take a look at the transactions within the ``FabCar`` smart contract.
Navigate to the ``chaincode/fabcar/javascript/lib`` subdirectory at the root of
``fabric-samples`` and open ``fabcar.js`` in your editor.

See how our smart contract is defined using the ``Contract`` class:

.. code:: bash

  class FabCar extends Contract {...

Within this class structure, you'll see that we have the following
transactions defined: ``initLedger``, ``queryCar``, ``queryAllCars``,
``createCar``, and ``changeCarOwner``. For example:


.. code:: bash

  async queryCar(ctx, carNumber) {...}
  async queryAllCars(ctx) {...}

Let's take a closer look at the ``queryAllCars`` transaction to see how it
interacts with the ledger.

.. code:: bash

  async queryAllCars(ctx) {

    const startKey = 'CAR0';
    const endKey = 'CAR999';

    const iterator = await ctx.stub.getStateByRange(startKey, endKey);


This code defines the range of cars that ``queryAllCars`` will retrieve from the
ledger. Every car between ``CAR0`` and ``CAR999`` -- 1,000 cars in all, assuming
every key has been tagged properly -- will be returned by the query. The
remainder of the code iterates through the query results and packages them into
JSON for the application.

Below is a representation of how an application would call different
transactions in a smart contract. Each transaction uses a broad set of APIs such
as ``getStateByRange`` to interact with the ledger. You can read more about
these APIs in `detail
<https://fabric-shim.github.io/master/index.html?redirect=true>`_.

.. image:: images/RunningtheSample.png

We can see our ``queryAllCars`` transaction, and another called ``createCar``.
We will use this later in the tutorial to update the ledger, and add a new block
to the blockchain.

But first, go back to the ``query`` program and change the
``evaluateTransaction`` request to query ``CAR4``. The ``query`` program should
now look like this:

.. code:: bash

  const result = await contract.evaluateTransaction('queryCar', 'CAR4');

Save the program and navigate back to your ``fabcar/javascript`` directory.
Now run the ``query`` program again:

.. code:: bash

  node query.js

You should see the following:

.. code:: json

  Wallet path: ...fabric-samples/fabcar/javascript/wallet
  Transaction has been evaluated, result is:
  {"colour":"black","make":"Tesla","model":"S","owner":"Adriana"}

If you go back and look at the result from when the transaction was
``queryAllCars``, you can see that ``CAR4`` was Adriana’s black Tesla model S,
which is the result that was returned here.

We can use the ``queryCar`` transaction to query against any car, using its
key (e.g. ``CAR0``) and get whatever make, model, color, and owner correspond to
that car.

Great. At this point you should be comfortable with the basic query transactions
in the smart contract and the handful of parameters in the query program.

Time to update the ledger...

Updating the ledger
-------------------

Now that we’ve done a few ledger queries and added a bit of code, we’re ready to
update the ledger. There are a lot of potential updates we could make, but
let's start by creating a **new** car.

From an application perspective, updating the ledger is simple. An application
submits a transaction to the blockchain network, and when it has been
validated and committed, the application receives a notification that
the transaction has been successful. Under the covers this involves the process
of **consensus** whereby the different components of the blockchain network work
together to ensure that every proposed update to the ledger is valid and
performed in an agreed and consistent order.

.. image:: tutorial/write_first_app.diagram.2.png

Above, you can see the major components that make this process work. As well as
the multiple peers which each host a copy of the ledger, and optionally a copy
of the smart contract, the network also contains an ordering service. The
ordering service coordinates transactions for a network; it creates blocks
containing transactions in a well-defined sequence originating from all the
different applications connected to the network.

Our first update to the ledger will create a new car. We have a separate program
called ``invoke.js`` that we will use to make updates to the ledger. Just as with
queries, use an editor to open the program and navigate to the code block where
we construct our transaction and submit it to the network:

.. code:: bash

  await contract.submitTransaction('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom');

See how the applications calls the smart contract transaction ``createCar`` to
create a black Honda Accord with an owner named Tom. We use ``CAR12`` as the
identifying key here, just to show that we don't need to use sequential keys.

Save it and run the program:

.. code:: bash

  node invoke.js

If the invoke is successful, you will see output like this:

.. code:: bash

  Wallet path: ...fabric-samples/fabcar/javascript/wallet
  2018-12-11T14:11:40.935Z - info: [TransactionEventHandler]: _strategySuccess: strategy success for transaction "9076cd4279a71ecf99665aed0ed3590a25bba040fa6b4dd6d010f42bb26ff5d1"
  Transaction has been submitted

Notice how the ``invoke`` application interacted with the blockchain network
using the ``submitTransaction`` API, rather than ``evaluateTransaction``.

.. code:: bash

  await contract.submitTransaction('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom');

``submitTransaction`` is much more sophisticated than ``evaluateTransaction``.
Rather than interacting with a single peer, the SDK will send the
``submitTransaction`` proposal to every required organization's peer in the
blockchain network. Each of these peers will execute the requested smart
contract using this proposal, to generate a transaction response which it signs
and returns to the SDK. The SDK collects all the signed transaction responses
into a single transaction, which it then sends to the orderer. The orderer
collects and sequences transactions from every application into a block of
transactions. It then distributes these blocks to every peer in the network,
where every transaction is validated and committed. Finally, the SDK is
notified, allowing it to return control to the application.

.. note:: ``submitTransaction`` also includes a listener that checks to make
          sure the transaction has been validated and committed to the ledger.
          Applications should either utilize a commit listener, or
          leverage an API like ``submitTransaction`` that does this for you.
          Without doing this, your transaction may not have been successfully
          orderered, validated, and committed to the ledger.

``submitTransaction`` does all this for the application! The process by which
the application, smart contract, peers and ordering service work together to
keep the ledger consistent across the network is called consensus, and it is
explained in detail in this `section <./peers/peers.html>`_.

To see that this transaction has been written to the ledger, go back to
``query.js`` and change the argument from ``CAR4`` to ``CAR12``.

In other words, change this:

.. code:: bash

  const result = await contract.evaluateTransaction('queryCar', 'CAR4');

To this:

.. code:: bash

  const result = await contract.evaluateTransaction('queryCar', 'CAR12');

Save once again, then query:

.. code:: bash

  node query.js

Which should return this:

.. code:: bash

  Wallet path: ...fabric-samples/fabcar/javascript/wallet
  Transaction has been evaluated, result is:
  {"colour":"Black","make":"Honda","model":"Accord","owner":"Tom"}

Congratulations. You’ve created a car and verified that its recorded on the
ledger!

So now that we’ve done that, let’s say that Tom is feeling generous and he
wants to give his Honda Accord to someone named Dave.

To do this, go back to ``invoke.js`` and change the smart contract transaction
from ``createCar`` to ``changeCarOwner`` with a corresponding change in input
arguments:

.. code:: bash

  await contract.submitTransaction('changeCarOwner', 'CAR12', 'Dave');

The first argument --- ``CAR12`` --- identifies the car that will be changing
owners. The second argument --- ``Dave`` --- defines the new owner of the car.

Save and execute the program again:

.. code:: bash

  node invoke.js

Now let’s query the ledger again and ensure that Dave is now associated with the
``CAR12`` key:

.. code:: bash

  node query.js

It should return this result:

.. code:: bash

   Wallet path: ...fabric-samples/fabcar/javascript/wallet
   Transaction has been evaluated, result is:
   {"colour":"Black","make":"Honda","model":"Accord","owner":"Dave"}

The ownership of ``CAR12`` has been changed from Tom to Dave.

.. note:: In a real world application the smart contract would likely have some
          access control logic. For example, only certain authorized users may
          create new cars, and only the car owner may transfer the car to
          somebody else.

Summary
-------

Now that we’ve done a few queries and a few updates, you should have a pretty
good sense of how applications interact with a blockchain network using a smart
contract to query or update the ledger. You’ve seen the basics of the roles
smart contracts, APIs, and the SDK play in queries and updates and you should
have a feel for how different kinds of applications could be used to perform
other business tasks and operations.

Additional resources
--------------------

As we said in the introduction, we have a whole section on
:doc:`developapps/developing_applications` that includes in-depth information on
smart contracts, process and data design, a tutorial using a more in-depth
Commercial Paper `tutorial <./tutorial/commercial_paper.html>`_ and a large
amount of other material relating to the development of applications.

.. Licensed under Creative Commons Attribution 4.0 International License
   https://creativecommons.org/licenses/by/4.0/
=======
Writing Your First Application
==============================

.. note:: If you're not yet familiar with the fundamental architecture of a
          Fabric network, you may want to visit the :doc:`blockchain` and
          :doc:`build_network` documentation prior to continuing.

In this section we'll be looking at a handful of sample programs to see how Fabric
apps work. These apps (and the smart contract they use) -- collectively known as
``fabcar`` -- provide a broad demonstration of Fabric functionality.  Notably, we
will show the process for interacting with a Certificate Authority and generating
enrollment certificates, after which we will leverage these generated identities
(user objects) to query and update a ledger.

We’ll go through three principle steps:

  **1. Setting up a development environment.** Our application needs a network to
  interact with, so we'll download one stripped down to just the components we need
  for registration/enrollment, queries and updates:

  .. image:: images/AppConceptsOverview.png

  **2. Learning the parameters of the sample smart contract our app will use.** Our
  smart contract contains various functions that allow us to interact with the ledger
  in different ways. We’ll go in and inspect that smart contract to learn about the
  functions our applications will be using.

  **3. Developing the applications to be able to query and update assets on the ledger.**
  We'll get into the app code itself (our apps have been written in Javascript) and
  manually manipulate the variables to run different kinds of queries and updates.

After completing this tutorial you should have a basic understanding of how
an application is programmed in conjunction with a smart contract to interact
with the ledger (i.e. the peer) on a Fabric network.

Setting up your Dev Environment
-------------------------------

First thing, let's download the Fabric images and the accompanying artifacts for the network
and applications...

Visit the :doc:`prereqs` page and ensure you have the necessary dependencies
installed on your machine.

Next, visit the :doc:`samples` page and follow the provided instructions.  Return to
this tutorial once you have cloned the ``fabric-samples`` repository, and downloaded
the latest stable Fabric images and available utilities.

At this point everything should be installed.  Navigate to the ``fabcar`` subdirectory
within your ``fabric-samples`` repository and take a look at what's inside:

.. code:: bash

  cd fabric-samples/fabcar  && ls

You should see the following:

.. code:: bash

     enrollAdmin.js	invoke.js	package.json	query.js	registerUser.js	startFabric.sh

Before starting we also need to do a little housekeeping.  Run the following command to
kill any stale or active containers:

.. code:: bash

  docker rm -f $(docker ps -aq)

Clear any cached networks:

.. code:: bash

  # Press 'y' when prompted by the command

  docker network prune

And lastly if you've already run through this tutorial, you'll also want to delete the
underlying chaincode image for the ``fabcar`` smart contract.  If you're a user going through
this content for the first time, then you won't have this chaincode image on your system:

.. code:: bash

  docker rmi dev-peer0.org1.example.com-fabcar-1.0-5c906e402ed29f20260ae42283216aa75549c571e2e380f3615826365d8269ba

Install the clients & launch the network
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. note:: The following instructions require you to be in the ``fabcar`` subdirectory
          within your local clone of the ``fabric-samples`` repo.  Remain at the
          root of this subdirectory for the remainder of this tutorial.

Run the following command to install the Fabric dependencies for the applications.
We are concerned with ``fabric-ca-client`` which will allow our app(s) to communicate
with the CA server and retrieve identity material, and with ``fabric-client`` which
allows us to load the identity material and talk to the peers and ordering service.

.. code:: bash

  npm install

Launch your network using the ``startFabric.sh`` shell script.  This command
will spin up our various Fabric entities and launch a smart contract container for
chaincode written in Golang:

.. code:: bash

  ./startFabric.sh

Alright, now that you’ve got a sample network and some code, let’s take a
look at how the different pieces fit together.

How Applications Interact with the Network
------------------------------------------

For a more in-depth look at the components in our ``fabcar`` network (and how
they're deployed) as well as how applications interact with those components
on more of a granular level, see :doc:`understand_fabcar_network`.

Developers more interested in seeing what applications **do** -- as well as
looking at the code itself to see how an application is constructed -- should
continue. For now, the most important thing to know is that applications use
a software development kit (SDK) to access the **APIs** that permit queries and
updates to the ledger.

Enrolling the Admin User
------------------------

.. note:: The following two sections involve communication with the Certificate
          Authority.  You may find it useful to stream the CA logs when running
          the upcoming programs.

To stream your CA logs, split your terminal or open a new shell and issue the following:

.. code:: bash

  docker logs -f ca.example.com

Now hop back to your terminal with the ``fabcar`` content...

When we launched our network, an admin user - ``admin`` - was registered with our
Certificate Authority.  Now we need to send an enroll call to the CA server and
retrieve the enrollment certificate (eCert) for this user.  We won't delve into enrollment
details here, but suffice it to say that the SDK and by extension our applications
need this cert in order to form a user object for the admin.  We will then use this admin
object to subsequently register and enroll a new user.  Send the admin enroll call to the CA
server:

.. code:: bash

  node enrollAdmin.js

This program will invoke a certificate signing request (CSR) and ultimately output
an eCert and key material into a newly created folder - ``hfc-key-store`` - at the
root of this project.  Our apps will then look to this location when they need to
create or load the identity objects for our various users.

Register and Enroll ``user1``
-----------------------------

With our newly generated admin eCert, we will now communicate with the CA server
once more to register and enroll a new user.  This user - ``user1`` - will be
the identity we use when querying and updating the ledger.  It's important to
note here that it is the ``admin`` identity that is issuing the registration and
enrollment calls for our new user (i.e. this user is acting in the role of a registrar).
Send the register and enroll calls for ``user1``:

.. code:: bash

  node registerUser.js

Similar to the admin enrollment, this program invokes a CSR and outputs the keys
and eCert into the ``hfc-key-store`` subdirectory.  So now we have identity material for two
separate users - ``admin`` & ``user1``.  Time to interact with the ledger...

Querying the Ledger
-------------------

Queries are how you read data from the ledger. This data is stored as a series
of key/value pairs, and you can query for the value of a single key, multiple
keys, or -- if the ledger is written in a rich data storage format like JSON --
perform complex searches against it (looking for all assets that contain
certain keywords, for example).

This is a representation of how a query works:

.. image:: images/QueryingtheLedger.png

First, let's run our ``query.js`` program to return a listing of all the cars on
the ledger.  We will use our second identity - ``user1`` - as the signing entity
for this application.  The following line in our program specifies ``user1`` as
the signer:

.. code:: bash

  fabric_client.getUserContext('user1', true);

Recall that the ``user1`` enrollment material has already been placed into our
``hfc-key-store`` subdirectory, so we simply need to tell our application to grab that identity.
With the user object defined, we can now proceed with reading from the ledger.
A function that will query all the cars, ``queryAllCars``, is
pre-loaded in the app, so we can simply run the program as is:

.. code:: bash

  node query.js

It should return something like this:

.. code:: json

  Query result count =  1
  Response is  [{"Key":"CAR0", "Record":{"colour":"blue","make":"Toyota","model":"Prius","owner":"Tomoko"}},
  {"Key":"CAR1",   "Record":{"colour":"red","make":"Ford","model":"Mustang","owner":"Brad"}},
  {"Key":"CAR2", "Record":{"colour":"green","make":"Hyundai","model":"Tucson","owner":"Jin Soo"}},
  {"Key":"CAR3", "Record":{"colour":"yellow","make":"Volkswagen","model":"Passat","owner":"Max"}},
  {"Key":"CAR4", "Record":{"colour":"black","make":"Tesla","model":"S","owner":"Adriana"}},
  {"Key":"CAR5", "Record":{"colour":"purple","make":"Peugeot","model":"205","owner":"Michel"}},
  {"Key":"CAR6", "Record":{"colour":"white","make":"Chery","model":"S22L","owner":"Aarav"}},
  {"Key":"CAR7", "Record":{"colour":"violet","make":"Fiat","model":"Punto","owner":"Pari"}},
  {"Key":"CAR8", "Record":{"colour":"indigo","make":"Tata","model":"Nano","owner":"Valeria"}},
  {"Key":"CAR9", "Record":{"colour":"brown","make":"Holden","model":"Barina","owner":"Shotaro"}}]

These are the 10 cars. A black Tesla Model S owned by Adriana, a red Ford Mustang
owned by Brad, a violet Fiat Punto owned by Pari, and so on. The ledger is
key/value based and in our implementation the key is ``CAR0`` through ``CAR9``.
This will become particularly important in a moment.

Let's take a closer look at this program.  Use an editor (e.g. atom or visual studio)
and open ``query.js``.

The initial section of the application defines certain variables such as
channel name, cert store location and network endpoints. In our sample app, these
variables have been baked-in, but in a real app these variables would have to
be specified by the app dev.

.. code:: bash

  var channel = fabric_client.newChannel('mychannel');
  var peer = fabric_client.newPeer('grpc://localhost:7051');
  channel.addPeer(peer);

  var member_user = null;
  var store_path = path.join(__dirname, 'hfc-key-store');
  console.log('Store path:'+store_path);
  var tx_id = null;

This is the chunk where we construct our query:

.. code:: bash

  // queryCar chaincode function - requires 1 argument, ex: args: ['CAR4'],
  // queryAllCars chaincode function - requires no arguments , ex: args: [''],
  const request = {
    //targets : --- letting this default to the peers assigned to the channel
    chaincodeId: 'fabcar',
    fcn: 'queryAllCars',
    args: ['']
  };

When the application ran, it invoked the ``fabcar`` chaincode on the peer, ran the
``queryAllCars`` function within it, and passed no arguments to it.

To take a look at the available functions within our smart contract, navigate
to the ``chaincode/fabcar/go`` subdirectory at the root of ``fabric-samples`` and open
``fabcar.go`` in your editor.

.. note:: These same functions are defined within the Node.js version of the
          ``fabcar`` chaincode.

You'll see that we have the following functions available to call: ``initLedger``,
``queryCar``, ``queryAllCars``, ``createCar``, and ``changeCarOwner``.

Let's take a closer look at the ``queryAllCars`` function to see how it
interacts with the ledger.

.. code:: bash

  func (s *SmartContract) queryAllCars(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "CAR0"
	endKey := "CAR999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)

This defines the range of ``queryAllCars``. Every car between ``CAR0`` and
``CAR999`` -- 1,000 cars in all, assuming every key has been tagged properly
-- will be returned by the query.

Below is a representation of how an app would call different functions in
chaincode.  Each function must be coded against an available API in the chaincode
shim interface, which in turn allows the smart contract container to properly
interface with the peer ledger.

.. image:: images/RunningtheSample.png

We can see our ``queryAllCars`` function, as well as one called ``createCar``,
that will allow us to update the ledger and ultimately append a new block to
the chain in a moment.

But first, go back to the ``query.js`` program and edit the constructor request
to query ``CAR4``. We do this by changing the function in ``query.js`` from
``queryAllCars`` to ``queryCar`` and passing ``CAR4`` as the specific key.

The ``query.js`` program should now look like this:

.. code:: bash

  const request = {
    //targets : --- letting this default to the peers assigned to the channel
    chaincodeId: 'fabcar',
    fcn: 'queryCar',
    args: ['CAR4']
  };

Save the program and navigate back to your ``fabcar`` directory.  Now run the
program again:

.. code:: bash

  node query.js

You should see the following:

.. code:: json

  {"colour":"black","make":"Tesla","model":"S","owner":"Adriana"}

If you go back and look at the result from when we queried every car before,
you can see that ``CAR4`` was Adriana’s black Tesla model S, which is the result
that was returned here.

Using the ``queryCar`` function, we can query against any key (e.g. ``CAR0``)
and get whatever make, model, color, and owner correspond to that car.

Great. At this point you should be comfortable with the basic query functions
in the smart contract and the handful of parameters in the query program.
Time to update the ledger...

Updating the Ledger
-------------------

Now that we’ve done a few ledger queries and added a bit of code, we’re ready to
update the ledger. There are a lot of potential updates we could make, but
let's start by creating a car.

Below we can see how this process works. An update is proposed, endorsed,
then returned to the application, which in turn sends it to be ordered and
written to every peer's ledger:

.. image:: images/UpdatingtheLedger.png

Our first update to the ledger will be to create a new car.  We have a separate
Javascript program -- ``invoke.js`` -- that we will use to make updates. Just
as with queries, use an editor to open the program and navigate to the
code block where we construct our invocation:

.. code:: bash

  // createCar chaincode function - requires 5 args, ex: args: ['CAR12', 'Honda', 'Accord', 'Black', 'Tom'],
  // changeCarOwner chaincode function - requires 2 args , ex: args: ['CAR10', 'Barry'],
  // must send the proposal to endorsing peers
  var request = {
    //targets: let default to the peer assigned to the client
    chaincodeId: 'fabcar',
    fcn: '',
    args: [''],
    chainId: 'mychannel',
    txId: tx_id
  };

You'll see that we can call one of two functions - ``createCar`` or
``changeCarOwner``. First, let’s create a red Chevy Volt and give it to an
owner named Nick. We're up to ``CAR9`` on our ledger, so we'll use ``CAR10``
as the identifying key here. Edit this code block to look like this:

.. code:: bash

  var request = {
    //targets: let default to the peer assigned to the client
    chaincodeId: 'fabcar',
    fcn: 'createCar',
    args: ['CAR10', 'Chevy', 'Volt', 'Red', 'Nick'],
    chainId: 'mychannel',
    txId: tx_id
  };

Save it and run the program:

.. code:: bash

   node invoke.js

There will be some output in the terminal about ``ProposalResponse`` and
promises. However, all we're concerned with is this message:

.. code:: bash

   The transaction has been committed on peer localhost:7053

To see that this transaction has been written, go back to ``query.js`` and
change the argument from ``CAR4`` to ``CAR10``.

In other words, change this:

.. code:: bash

  const request = {
    //targets : --- letting this default to the peers assigned to the channel
    chaincodeId: 'fabcar',
    fcn: 'queryCar',
    args: ['CAR4']
  };

To this:

.. code:: bash

  const request = {
    //targets : --- letting this default to the peers assigned to the channel
    chaincodeId: 'fabcar',
    fcn: 'queryCar',
    args: ['CAR10']
  };

Save once again, then query:

.. code:: bash

  node query.js

Which should return this:

.. code:: bash

   Response is  {"colour":"Red","make":"Chevy","model":"Volt","owner":"Nick"}

Congratulations. You’ve created a car!

So now that we’ve done that, let’s say that Nick is feeling generous and he
wants to give his Chevy Volt to someone named Dave.

To do this go back to ``invoke.js`` and change the function from ``createCar``
to ``changeCarOwner`` and input the arguments like this:

.. code:: bash

  var request = {
    //targets: let default to the peer assigned to the client
    chaincodeId: 'fabcar',
    fcn: 'changeCarOwner',
    args: ['CAR10', 'Dave'],
    chainId: 'mychannel',
    txId: tx_id
  };

The first argument -- ``CAR10`` -- reflects the car that will be changing
owners. The second argument -- ``Dave`` -- defines the new owner of the car.

Save and execute the program again:

.. code:: bash

  node invoke.js

Now let’s query the ledger again and ensure that Dave is now associated with the
``CAR10`` key:

.. code:: bash

  node query.js

It should return this result:

.. code:: bash

   Response is  {"colour":"Red","make":"Chevy","model":"Volt","owner":"Dave"}

The ownership of ``CAR10`` has been changed from Nick to Dave.

.. note:: In a real world application the chaincode would likely have some access
          control logic. For example, only certain authorized users may create
          new cars, and only the car owner may transfer the car to somebody else.

Summary
-------

Now that we’ve done a few queries and a few updates, you should have a pretty
good sense of how applications interact with the network. You’ve seen the basics
of the roles smart contracts, APIs, and the SDK play in queries and updates and
you should have a feel for how different kinds of applications could be used to
perform other business tasks and operations.

In subsequent documents we’ll learn how to actually **write** a smart contract
and how some of these more low level application functions can be leveraged
(especially relating to identity and membership services).

Additional Resources
--------------------

The `Hyperledger Fabric Node SDK repo <https://github.com/hyperledger/fabric-sdk-node>`__
is an excellent resource for deeper documentation and sample code.  You can also consult
the Fabric community and component experts on `Hyperledger Rocket Chat <https://chat.hyperledger.org/home>`__.

.. Licensed under Creative Commons Attribution 4.0 International License
   https://creativecommons.org/licenses/by/4.0/
>>>>>>> release-1.0
