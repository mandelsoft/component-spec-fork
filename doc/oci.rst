Component Content in an OCI Registry
====================================

Component descriptors and local content delivered with a component are stored
in an OCI registry as OCI artefact. The component descriptor may describe
references to other artefacts in any repository what type is supported
by an access type, regardless of the type of the artefact. There are dedicated
access types to refer to artefacts stored in an OCI registry.

Storing Components
------------------

The content of a component is stored in the OCI registry as OCI artefact of the media type ``application/vnd.gardener.cloud.cnudie.component.config.v1+json``. Hereby the content consists at least of the component descriptor itself, which is basically stored as blob described as first layer in the OCI manifest.
The media type is used as type for the config blob described in the OCI manifest. The blob contains a config json describing the digest of the component descriptor blob (first layer). (describing the blob digest would be sufficient, but so far we store this blob always as first layer)

The component descriptor blob uses the media type ``application/vnd.gardener.cloud.cnudie.component-descriptor.v2+yaml``. It is stored directly as blob.

Together with the component descriptor a set of artefacts local to the component may be delivered. Those component local artefact blobs are stored as additional layers of the appropriate media type. To describe the access to such artefacts local to the component the component descriptor supports a dedicated access type ``LocalBlob`` with the blob digest as reference name.

Creating Components during the Build
------------------------------------

Typically versions for components are provided by builds. 
To support this scenario there is a contract between a build and a
component uploader tool.

Te task of the build is then just to provide such a contract file
containing the generated component descriptor version with approriate
artefact and component references together with referenced blob artefacts.
To describe references to those local artefact blobs the 
well-known access type ``LocalBlob`` is used.

A build may provide information for multiple components. For every provided
component version a component archive with the component descriptor and the 
contained local blobs has to be provided. Those archives are then bundled
in another fingle archive. THis archive is then the contract to the upload
tool. It is possible to call this tool directly in the build or the build system
defined a way to pass this archive to the build ecosystem which then carries
the repository credentials and handles the upload.

The same contract can be used for transporting repository content
via storage media between two environments. Therefore the transport tool
just needs to provide a file based frontend and backend additional to
the regular repository to repository endpoints. This leads to the
specification of a general CNUDIE Transport Format (CTF).

The same contract and tooling can then be used for the build uploads and
regular transports.

.. image:: Build.png

CNUDIE Transport Format
-----------------------

The transport format is based on tar archives (.tar or .tgz). There is one
archive per OCI artefact that should be included in the transport.

There will be two similar file system structures, one to describe component
artefact and one for describing regular standard OCI artefacts.

The final transport file is then just an archive (tar, tgz) containing the
artefact archives.

Component Archive
.................

So far there is only a specification for the format for a component artefact.

.. code-block::

  ├── component-descriptor.yaml
  └── blobs
      ├── blob1
      ├──  ...
      └── blobn

All the contained blobs (may use any name here) have to be referenced by the
local access type in the component descriptor.

The component descriptor must be the first entry in the tar archive to
support streaming.


.. _oci-artefact:

OCI Artefact
............

A similar format for standard OCI Artefacts (including OCI Images as special case)
could look like this

.. code-block::

  ├── artefact-descriptor.yaml
  ├── oci-manifest.json
  └── blobs
      ├── blob1
      ├──  ...
      └── blobn

Like the component descriptor the additional file ``artefact-descriptor.yaml``
describes the artefact name and version. The other files are just taken
from the OCI artefact api. The blob names should be the correct OCI
digests used in the ``oci-manifest.json``.

The artefact descriptor must be the first entry in the tar archive to
support streaming, followed by the OCI manifest.

OCI Blob
............

Dedicated blobs stored as OCI artefact could be represented in a simplified manner (or line OCI Artefacts).

.. code-block::

  ├── artefact-descriptor.yaml
  └── blobs
      └── blob

Like for the OCI artefact the additional file ``artefact-descriptor.yaml``
describes the artefact name and version. Additionally it must describe
the mime type of the blob. The ``oci-manifest.json`` for this representation is then automatically be created during the upload.

The artefact descriptor must be the first entry in the tar archive to
support streaming, followed by the blob.


OCI related Access Types
------------------------

To be able to describe access to artefacts stored in an OCI repository the 
component descriptor supports several dedicated access types.
There are two ways to reference OCI stored artefacts:
- using the fully qualified access url to address artefacts stored in repositories different from the one hosting the using component descriptor.
- using repository local names used to refer artefacts in the same repository. Supporting such local references enables the usage of technical replication tools for copying complete content of repositories without knowing about the component model and the structure of component descriptor files.

Additionally there are two kinds of artefacts that have to be addressed in an
OCI registry:
- Direct blob artefacts (using the digest based blob access for a registry)
- OCI artefacts consisting of multiple blobs described by an OCI manifest.
- Blobs described as OCI artefact. Potentially used to decouple component local blobs from their containing content during a transport step

For direct blobs there is also the possibility to store such artefacts directly as part of the content of a component.

In summary there are therefore seven access combinations, that are described by dedicated access types:

``OCIBlob``
  A fully qualified URL for accessing a blob using the OCI blob api using a repository URL, a blob digest and a repository path.

``OCIArtefact``
  A fully qualified URL for accessing an OCI aretfact using the OCI artefact api using a repository URL, a manifest digest/version and a respository path.

``OCIBlobArtefact``
  A blob explicitly stored via an OCI artefact with an own identity independent of a dedicated component.
  It is stored as layer with the appropriate media type. The config.json so far has no content. It can be used to store any blob directly as blob with exposing the artefact as part of the resource blob.

``RepositoryLocalOCIBlob``
  A path in the local repository (the artefact name) and the blob digest used to access the blob via the OCI blob api in the repository hosting the artefact reference.

``RepositoryLocalOCIArtefact``
  A path in the local repository (the artefact name) and the manifest digest used to access the blob via the OCI artefact api in the repository hosting the artefact reference.

``RepositoryLocalOCIBlobArtefact``
  A path in the local repository (the artefact name) denoting an OCI artefact used to represent a single blob stored as sole layer.

``LocalBlob``
  The digest of the blob using the OCI blob api to access the blob as blob nested to the actual component in the repository hosting the component descriptor.


.. image:: Blob.png

Resource Access for Component Resource
--------------------------------------

A library evaluating a component descriptor MUST provide functions to access
the resources described by the component descriptor. Therefore it has to implement the various access types.

Hereby the blob types ( ``LocalBlob``, ``RepositoryLocalOCIBlob``, ``OCIBlob``, ``RepositoryLocalOCIBlobArtefact`` and ``OCIBlobArtefact``) directly return the denoted blobs (stored as layer or directly as blob). The OCI artefact types (``OCIArtefact`` and ``RepositoryLocalOCIArtefact``) return the content of the complete OCI artefact as tar blob as described in section :ref:`oci-artefact` used for the transport format.

