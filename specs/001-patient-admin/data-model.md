# Data Model: Patient Administration Sub-System

**Feature**: Patient Administration Sub-System  
**Date**: 2025-12-20

This document lists primary entities, fields, relationships, validation rules, and state transitions.

## Entities

### Patient
- id: UUID (primary)
- mrn: string (medical record number, unique within organization)
- primary_person_id: UUID (ref → Person)
- dob: date
- gender: enum (M/F/Other/Unknown)
- contact_preferences: JSON (email, sms, phone, consent flags)
- consents: [Consent] (references to consent records)
- created_at, updated_at

Validation: mrn is required; dob must be in the past; consent flags validated as booleans.

### Person
- id: UUID
- given_name, family_name
- emails: list (validated email objects)
- phones: list (validated phone objects)
- addresses: list
- type: enum (patient, provider, contact)
- created_at, updated_at

### Provider
- id: UUID
- person_id: UUID (ref → Person)
- credentials: list (licenses, NPI)
- specialties: list
- schedules: list (availability blocks)
- active: boolean
- created_at, updated_at

### Location / ServiceDeliveryLocation / Place
- id: UUID
- name: string
- address: structured address
- type: enum (clinic, room, virtual)
- active: boolean

### Organization
- id: UUID
- name, tax_id
- contact info

### Resource
- id: UUID
- type: enum (room, equipment, staff)
- attributes: JSON

### Schedule
- id: UUID
- owner_type: enum (provider, location, resource)
- owner_id: UUID
- recurring_rules: RRULE or availability blocks
- start_time, end_time (availability window)
- timezone
- created_at, updated_at

### Slot
- id: UUID
- schedule_id: UUID (ref → Schedule)
- start_time, end_time
- status: enum (free, held, booked, unavailable)
- hold_token: nullable string (reservation token when held)
- created_at, updated_at

### Consent
- id: UUID
- patient_id: UUID (ref → Patient)
- purpose: string (e.g., appointment_reminder, marketing - marketing OUT OF SCOPE)
- channel: enum (email, sms, in_app)
- granted_at: timestamp
- revoked_at: timestamp (nullable)
- source: enum (ui, api, scim)
- granted_by: string (user or system)
- metadata: JSON (e.g., ip_address, user_agent)
- audit_trail: immutable history entries
- Retention: consent records and revocation history **MUST** be retained per SC-SEC-004 (6 years).

### Appointment / Booking
- id: UUID
- patient_id: UUID
- participants: [{role: clinician|patient|staff, actor_id}]
- provider_id: UUID
- location_id: UUID
- resources: [Resource refs]
- start_time, end_time
- status: enum (requested, pending, booked, confirmed, cancelled, completed, noshow)
- created_by, updated_by
- change_history: audit trail (timestamps, actor, changes)

Validation: start_time < end_time; resources available; provider available in schedule.

## Relationships
- Patient → Person (1:1 primary person; 1:many for contacts)
- Provider → Person (1:1)
- Organization → Provider (1:many)
- Location → Organization (1:many)
- Appointment → Patient, Provider, Location, Resource(s)

## State Transitions (Appointment)
- requested → pending (system validates availability)
- pending → booked (resources and provider reserved)
- booked → confirmed (stakeholder confirmed)
- booked/confirmed → cancelled (cancellation policy may apply)
- confirmed → completed (after scheduled time and confirmation)
- any state → noshow (if not attended)

Events and invariants:
- On transition → audit event recorded (actor, time, state delta)
- Idempotency keys for booking operations must ensure no double-booking from retries

## Index & Query Strategies
- Index Patient.mrn, Patient.primary_person_id
- Index Appointment by provider_id + start_time (for fast slot queries)
- Use materialized views or denormalized tables for clinician appointment lists if needed for performance

## Data Retention & Archival
- Appointment and audit records retained per policy (6 years). Older records archived to S3 (encrypted) per lifecycle rules.

## Notes
- Entities are intentionally normalized to enforce referential integrity and support consistent audit trails.  
- Schema will be emitted as migration scripts in Phase 1 implementation.  

