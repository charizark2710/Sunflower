export async function GET(
  request: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;

    const response = await fetch(
      `https://sunflower-rdips.tech/api/performances/${id}`,
      {
        headers: {
          'X-API-Key': 'a706913c38b555a889218175',
        },
      }
    );

    if (!response.ok) {
      return new Response(
        JSON.stringify({ error: 'Failed to fetch performances data' }),
        { status: response.status, headers: { 'Content-Type': 'application/json' } }
      );
    }

    const data = await response.json();
    return new Response(JSON.stringify(data), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    });
  } catch (error) {
    console.error('Error fetching performances:', error);
    return new Response(
      JSON.stringify({ error: 'Internal server error' }),
      { status: 500, headers: { 'Content-Type': 'application/json' } }
    );
  }
}
