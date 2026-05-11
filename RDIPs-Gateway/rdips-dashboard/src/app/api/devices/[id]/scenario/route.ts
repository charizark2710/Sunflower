import type { NextRequest } from 'next/server';

const API_BASE_URL = 'https://sunflower-rdips.tech/api';
const API_KEY = 'a706913c38b555a889218175';

type ScenarioType = 'UNDERESTIMATE' | 'OVERESTIMATE' | 'NORMAL';

interface ScenarioRequestBody {
  type: ScenarioType;
  /** Optional override for the time window length, in hours. Defaults to 4h. */
  durationHours?: number;
  /** Optional override for the interval between generated records, in minutes. Defaults to 60m. */
  intervalMinutes?: number;
}

interface PayloadRecord {
  createdAt: string;
  maxCapacity: string;
  inCapacity: string;
  outCapacity: string;
  capacity: string;
  uom: string;
}

interface DeviceDetailResponse {
  httpCode: number;
  data: {
    id: string;
    name: string;
    performanceID: string;
    performance?: {
      id: string;
      document_name: string;
    };
    [key: string]: unknown;
  };
  message: string;
  Error: unknown;
}

/**
 * Format a Date instance as "D/M/YYYY HH:MM:SS" (e.g. "9/5/2026 10:00:00").
 */
function formatCreatedAt(date: Date): string {
  const day = date.getDate();
  const month = date.getMonth() + 1;
  const year = date.getFullYear();
  const hh = String(date.getHours()).padStart(2, '0');
  const mm = String(date.getMinutes()).padStart(2, '0');
  const ss = String(date.getSeconds()).padStart(2, '0');
  return `${day}/${month}/${year} ${hh}:${mm}:${ss}`;
}

/**
 * Return a random float in [min, max], rounded to `decimals` places.
 */
function randomInRange(min: number, max: number, decimals = 2): number {
  const value = Math.random() * (max - min) + min;
  const factor = 10 ** decimals;
  return Math.round(value * factor) / factor;
}

/**
 * Build the payload array for the requested scenario type.
 *
 * Records are generated across `durationHours` (default 4h), spaced
 * `intervalMinutes` apart (default 60m). The relationship between
 * maxCapacity and capacity encodes the scenario:
 *
 *   UNDERESTIMATE → maxCapacity > capacity  (Status Warning)
 *   OVERESTIMATE  → capacity > maxCapacity  (Status Warning)
 *   NORMAL        → maxCapacity == capacity (Status Normal)
 *
 * Within each scenario, the actual numeric values are randomized so
 * that consecutive runs produce different-looking data.
 */
function buildPayload(
  scenario: ScenarioType,
  durationHours: number,
  intervalMinutes: number
): PayloadRecord[] {
  const records: PayloadRecord[] = [];
  const now = new Date();
  const totalMinutes = durationHours * 60;
  const stepCount = Math.max(1, Math.floor(totalMinutes / intervalMinutes));

  for (let i = 0; i <= stepCount; i++) {
    const createdAtDate = new Date(now.getTime() + i * intervalMinutes * 60 * 1000);

    // Random base in [8, 14] so each record uses a different magnitude.
    const baseMax = randomInRange(8, 14, 2);

    let maxCapacity: number;
    let capacity: number;

    switch (scenario) {
      case 'UNDERESTIMATE': {
        // maxCapacity > capacity → warning (system under-estimated demand)
        maxCapacity = baseMax;
        // capacity is 30%–80% of maxCapacity → strictly less than maxCapacity
        capacity = randomInRange(maxCapacity * 0.3, maxCapacity * 0.8, 2);
        break;
      }
      case 'OVERESTIMATE': {
        // capacity > maxCapacity → warning (system over-estimated demand)
        maxCapacity = baseMax;
        // capacity is 120%–180% of maxCapacity → strictly greater than maxCapacity
        capacity = randomInRange(maxCapacity * 1.2, maxCapacity * 1.8, 2);
        break;
      }
      case 'NORMAL':
      default: {
        // maxCapacity == capacity → normal
        maxCapacity = baseMax;
        capacity = baseMax;
        break;
      }
    }

    // inCapacity / outCapacity are also randomized but kept proportional
    // to capacity so the numbers feel realistic.
    const inCapacity = randomInRange(capacity * 0.1, capacity * 0.4, 2);
    const outCapacity = randomInRange(capacity * 0.02, capacity * 0.15, 2);

    records.push({
      createdAt: formatCreatedAt(createdAtDate),
      maxCapacity: String(maxCapacity),
      inCapacity: String(inCapacity),
      outCapacity: String(outCapacity),
      capacity: String(capacity),
      uom: 'kw',
    });
  }

  return records;
}

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;

    // --- 1. Parse + validate request body --------------------------------
    let body: ScenarioRequestBody;
    try {
      body = (await request.json()) as ScenarioRequestBody;
    } catch {
      return jsonResponse({ error: 'Invalid JSON body' }, 400);
    }

    const allowed: ScenarioType[] = ['UNDERESTIMATE', 'OVERESTIMATE', 'NORMAL'];
    if (!body?.type || !allowed.includes(body.type)) {
      return jsonResponse(
        {
          error:
            'Invalid "type". Must be one of: UNDERESTIMATE, OVERESTIMATE, NORMAL.',
        },
        400
      );
    }

    const durationHours =
      typeof body.durationHours === 'number' && body.durationHours > 0
        ? body.durationHours
        : 4;
    const intervalMinutes =
      typeof body.intervalMinutes === 'number' && body.intervalMinutes > 0
        ? body.intervalMinutes
        : 60;

    // --- 2. GET device detail to obtain performanceID + device name ------
    const deviceRes = await fetch(
      `${API_BASE_URL}/devices/${id}?detail=true`,
      {
        method: 'GET',
        headers: {
          'X-API-Key': API_KEY,
        },
        cache: 'no-store',
      }
    );

    if (!deviceRes.ok) {
      const text = await deviceRes.text();
      return jsonResponse(
        {
          error: 'Failed to fetch device detail',
          upstreamStatus: deviceRes.status,
          upstreamBody: text,
        },
        deviceRes.status
      );
    }

    const deviceJson = (await deviceRes.json()) as DeviceDetailResponse;
    const deviceID = deviceJson?.data.id;
    const documentName = deviceJson?.data?.name;

    if (!deviceID) {
      return jsonResponse(
        { error: 'Device has no performanceID', device: deviceJson?.data },
        422
      );
    }

    // --- 3. Build payload + PUT to performances endpoint -----------------
    const payload = buildPayload(body.type, durationHours, intervalMinutes);

    const putRes = await fetch(
      `${API_BASE_URL}/performances/${deviceID}`,
      {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': API_KEY,
        },
        body: JSON.stringify({
          document_name: documentName,
          payload,
        }),
      }
    );

    const putText = await putRes.text();
    let putBody: unknown = putText;
    try {
      putBody = JSON.parse(putText);
    } catch {
      // keep raw text if not JSON
    }

    if (!putRes.ok) {
      return jsonResponse(
        {
          error: 'Failed to update performance',
          upstreamStatus: putRes.status,
          upstreamBody: putBody,
        },
        putRes.status
      );
    }

    return jsonResponse({
      ok: true,
      scenario: body.type,
      deviceId: id,
      deviceID,
      documentName,
      durationHours,
      intervalMinutes,
      recordCount: payload.length,
      payload,
      upstream: putBody,
    });
  } catch (err) {
    console.error('Error in POST /api/devices/[id]/scenario:', err);
    return jsonResponse(
      {
        error: 'Internal server error',
        detail: err instanceof Error ? err.message : String(err),
      },
      500
    );
  }
}
