using Microsoft.AspNetCore.Mvc;

namespace cartservice.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class CartController : ControllerBase
    {
        [HttpPost("add")]
        public IActionResult Add([FromBody] object payload)
        {
            // TODO: add real logic. This stub returns OK.
            return Ok(new { status = "added" });
        }
    }
}